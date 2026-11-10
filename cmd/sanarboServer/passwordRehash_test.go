package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/lao-tseu-is-alive/sanarbo/pkg/version"
)

// test passwords are already sha256 hex strings, as sent by the frontend
const (
	rehashTestUser     = "jdoe"
	rehashTestPassword = "4acf0b39d9c4766709a3689f553ac01ab550545ffa4544dfc0b2cea82fba02a3"
	// keep the tests fast: the production target is minBcryptCost
	rehashTestTargetCost = bcrypt.MinCost + 1
)

type fakePasswordHashStore struct {
	hashes    map[string]string
	getErr    error
	updateErr error
	updates   int
}

func (f *fakePasswordHashStore) GetPasswordHash(username string) (string, error) {
	if f.getErr != nil {
		return "", f.getErr
	}
	h, ok := f.hashes[username]
	if !ok {
		return "", errors.New("no rows in result set")
	}
	return h, nil
}

func (f *fakePasswordHashStore) UpdatePasswordHash(username, oldHash, newHash string) (int, error) {
	if f.updateErr != nil {
		return 0, f.updateErr
	}
	if f.hashes[username] != oldHash {
		return 0, nil
	}
	f.hashes[username] = newHash
	f.updates++
	return 1, nil
}

func mustHash(t *testing.T, password string, cost int) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	require.NoError(t, err)
	return string(h)
}

// fakeLoginHandler behaves like users.Service.LoginUser: it binds the body and checks the bcrypt hash.
func fakeLoginHandler(store *fakePasswordHashStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		login := new(users.UserLogin)
		if err := c.Bind(login); err != nil {
			return c.JSON(http.StatusBadRequest, "invalid user login or json format in request body")
		}
		stored, ok := store.hashes[login.Username]
		if !ok || bcrypt.CompareHashAndPassword([]byte(stored), []byte(login.PasswordHash)) != nil {
			return c.JSON(http.StatusUnauthorized, "wrong password")
		}
		return c.JSON(http.StatusOK, map[string]string{"token": "a.b.c"})
	}
}

func doLogin(t *testing.T, store *fakePasswordHashStore, body string) *httptest.ResponseRecorder {
	t.Helper()
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_rehash")
	require.NoError(t, err)
	e := echo.New()
	e.POST("/login", fakeLoginHandler(store), rehashPasswordOnLogin(store, l, rehashTestTargetCost))
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func loginBody(username, password string) string {
	return `{"username":"` + username + `","password_hash":"` + password + `"}`
}

func TestRehashPasswordOnLogin(t *testing.T) {
	t.Run("successful login with a cheap hash upgrades it", func(t *testing.T) {
		store := &fakePasswordHashStore{hashes: map[string]string{rehashTestUser: mustHash(t, rehashTestPassword, bcrypt.MinCost)}}

		rec := doLogin(t, store, loginBody(rehashTestUser, rehashTestPassword))

		assert.Equal(t, http.StatusOK, rec.Code, "the login handler must still receive the whole body")
		assert.Contains(t, rec.Body.String(), `"token"`)
		assert.Equal(t, 1, store.updates)
		cost, err := bcrypt.Cost([]byte(store.hashes[rehashTestUser]))
		require.NoError(t, err)
		assert.Equal(t, rehashTestTargetCost, cost)
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(store.hashes[rehashTestUser]), []byte(rehashTestPassword)),
			"the new hash must still accept the same password")
	})

	t.Run("failed login does not touch the hash", func(t *testing.T) {
		old := mustHash(t, rehashTestPassword, bcrypt.MinCost)
		store := &fakePasswordHashStore{hashes: map[string]string{rehashTestUser: old}}

		rec := doLogin(t, store, loginBody(rehashTestUser, strings.Repeat("0", 64)))

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, 0, store.updates)
		assert.Equal(t, old, store.hashes[rehashTestUser])
	})

	t.Run("hash already at target cost is left as is", func(t *testing.T) {
		old := mustHash(t, rehashTestPassword, rehashTestTargetCost)
		store := &fakePasswordHashStore{hashes: map[string]string{rehashTestUser: old}}

		rec := doLogin(t, store, loginBody(rehashTestUser, rehashTestPassword))

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 0, store.updates)
		assert.Equal(t, old, store.hashes[rehashTestUser])
	})

	t.Run("store errors never make the login fail", func(t *testing.T) {
		for name, store := range map[string]*fakePasswordHashStore{
			"read error":   {getErr: errors.New("db down")},
			"update error": {updateErr: errors.New("permission denied")},
		} {
			t.Run(name, func(t *testing.T) {
				// the fake login handler reads the map directly, so it still succeeds
				store.hashes = map[string]string{rehashTestUser: mustHash(t, rehashTestPassword, bcrypt.MinCost)}

				rec := doLogin(t, store, loginBody(rehashTestUser, rehashTestPassword))

				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, 0, store.updates)
			})
		}
	})

	t.Run("invalid body is left to the login handler", func(t *testing.T) {
		store := &fakePasswordHashStore{hashes: map[string]string{}}

		rec := doLogin(t, store, `{not json`)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, 0, store.updates)
	})
}

func TestUpgradePasswordHashDoesNotOverwriteAConcurrentChange(t *testing.T) {
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_rehash")
	require.NoError(t, err)
	old := mustHash(t, rehashTestPassword, bcrypt.MinCost)
	changed := mustHash(t, "another password", bcrypt.MinCost)
	store := &racingStore{fakePasswordHashStore: fakePasswordHashStore{hashes: map[string]string{rehashTestUser: old}}, changeTo: changed}

	upgradePasswordHash(store, l, rehashTestUser, rehashTestPassword, rehashTestTargetCost)

	assert.Equal(t, changed, store.hashes[rehashTestUser], "a password changed between read and update must be kept")
}

// racingStore simulates a password reset happening right after the hash was read.
type racingStore struct {
	fakePasswordHashStore
	changeTo string
}

func (r *racingStore) GetPasswordHash(username string) (string, error) {
	h, err := r.fakePasswordHashStore.GetPasswordHash(username)
	r.hashes[username] = r.changeTo
	return h, err
}

func TestMinBcryptCost(t *testing.T) {
	assert.GreaterOrEqual(t, minBcryptCost, 12, "OWASP recommends a bcrypt cost of at least 10, sanarbo uses 12")
	assert.LessOrEqual(t, minBcryptCost, bcrypt.MaxCost)
}

// TestPgPasswordHashStore checks the two SQL queries against the real go_user table without changing it:
// the update is always given a stale old hash, so it must match no row.
func TestPgPasswordHashStore(t *testing.T) {
	if _, ok := os.LookupEnv("DB_PASSWORD"); !ok {
		t.Skip("DB_PASSWORD is not set, no database to test against")
	}
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_rehash")
	require.NoError(t, err)
	dsn := config.GetPgDbDsnUrlFromEnvOrPanic(defaultDBIp, defaultDBPort, tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	db, err := database.GetInstance("pgx", dsn, runtime.NumCPU(), l)
	require.NoError(t, err)
	defer db.Close()
	store := pgPasswordHashStore{db: db}

	t.Run("unknown user", func(t *testing.T) {
		_, err := store.GetPasswordHash("sanarbo-test-no-such-user")
		assert.Error(t, err)
		rows, err := store.UpdatePasswordHash("sanarbo-test-no-such-user", "old", "new")
		assert.NoError(t, err)
		assert.Equal(t, 0, rows)
	})

	t.Run("existing user", func(t *testing.T) {
		username, err := db.GetQueryString("SELECT username FROM go_user ORDER BY id LIMIT 1;")
		if err != nil {
			t.Skipf("no user in go_user: %v", err)
		}
		stored, err := store.GetPasswordHash(username)
		require.NoError(t, err)
		_, err = bcrypt.Cost([]byte(stored))
		assert.NoError(t, err, "go_user.password_hash must be a bcrypt hash")

		rows, err := store.UpdatePasswordHash(username, "stale-hash", "new")
		assert.NoError(t, err)
		assert.Equal(t, 0, rows, "a stale old hash must not overwrite the stored one")
		after, err := store.GetPasswordHash(username)
		require.NoError(t, err)
		assert.Equal(t, stored, after)
	})
}
