package main

import (
	"bytes"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"golang.org/x/crypto/bcrypt"
)

// minBcryptCost is the cost below which a stored password hash is recomputed after a successful login.
// go-cloud-k8s-user-group hashes new, reset and admin passwords with bcrypt.MinCost (4), which is far too
// cheap to resist an offline attack if go_user leaks. The client sends sha256(password): the stored value
// is bcrypt(sha256(password)), so the hash can only be recomputed when that value is received at login.
const minBcryptCost = 12

const (
	passwordHashGet    = "SELECT password_hash FROM go_user WHERE username = $1;"
	passwordHashUpdate = "UPDATE go_user SET password_hash = $1 WHERE username = $2 AND password_hash = $3;"
)

type passwordHashStore interface {
	GetPasswordHash(username string) (string, error)
	// UpdatePasswordHash replaces oldHash by newHash, only if the stored hash is still oldHash.
	UpdatePasswordHash(username, oldHash, newHash string) (int, error)
}

type pgPasswordHashStore struct {
	db database.DB
}

func (s pgPasswordHashStore) GetPasswordHash(username string) (string, error) {
	return s.db.GetQueryString(passwordHashGet, username)
}

func (s pgPasswordHashStore) UpdatePasswordHash(username, oldHash, newHash string) (int, error) {
	return s.db.ExecActionQuery(passwordHashUpdate, newHash, username, oldHash)
}

// rehashPasswordOnLogin wraps the login handler: after a successful login, it recomputes the stored
// bcrypt hash with targetCost when its cost is lower. It never makes a login fail.
func rehashPasswordOnLogin(store passwordHashStore, l golog.MyLogger, targetCost int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			raw, err := io.ReadAll(req.Body)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "unable to read request body")
			}
			login := new(users.UserLogin)
			req.Body = io.NopCloser(bytes.NewReader(raw))
			bindErr := c.Bind(login)
			// the login handler reads the body again
			req.Body = io.NopCloser(bytes.NewReader(raw))

			if err := next(c); err != nil {
				return err
			}
			if bindErr != nil || c.Response().Status != http.StatusOK || login.Username == "" {
				return nil
			}
			upgradePasswordHash(store, l, login.Username, login.PasswordHash, targetCost)
			return nil
		}
	}
}

func upgradePasswordHash(store passwordHashStore, l golog.MyLogger, username, password string, targetCost int) {
	stored, err := store.GetPasswordHash(username)
	if err != nil {
		l.Warn("password rehash: cannot read hash of user %q: %v", username, err)
		return
	}
	cost, err := bcrypt.Cost([]byte(stored))
	if err != nil || cost >= targetCost {
		return
	}
	// the login handler already checked it, but never store a hash of something that does not match
	if bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)) != nil {
		return
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), targetCost)
	if err != nil {
		l.Warn("password rehash: cannot compute hash of user %q: %v", username, err)
		return
	}
	rows, err := store.UpdatePasswordHash(username, stored, string(newHash))
	if err != nil {
		l.Warn("password rehash: cannot store hash of user %q: %v", username, err)
		return
	}
	if rows == 1 {
		l.Info("password rehash: bcrypt cost of user %q raised from %d to %d", username, cost, targetCost)
	}
}
