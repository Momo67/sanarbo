package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/cristalhq/jwt/v4"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goserver"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"github.com/lao-tseu-is-alive/sanarbo/pkg/trees"
	"github.com/lao-tseu-is-alive/sanarbo/pkg/version"
)

const (
	defaultPort                = 9999
	defaultDBPort              = 5432
	defaultDBIp                = "127.0.0.1"
	defaultDBSslMode           = "prefer"
	defaultWebRootDir          = "sanarboFront/dist/"
	defaultSqlDbMigrationsPath = "db/migrations"
	defaultSecuredApi          = "/goapi/v1"
	defaultUsername            = "bill"
	defaultFakeStupidPass      = "board"
	charsetUTF8                = "charset=UTF-8"
	MIMEAppJSON                = "application/json"
	MIMEHtml                   = "text/html"
	MIMEAppJSONCharsetUTF8     = MIMEAppJSON + "; " + charsetUTF8
	MIMEHtmlCharsetUTF8        = MIMEHtml + "; " + charsetUTF8
	HeaderContentType          = "Content-Type"
)

// content holds our static web server content.
//
//go:embed sanarboFront/dist/*
var content embed.FS

// sqlMigrations holds our db migrations sql files using https://github.com/golang-migrate/migrate
// in the line above you SHOULD have the same path  as const defaultSqlDbMigrationsPath
//
//go:embed db/migrations/*.sql
var sqlMigrations embed.FS

var dbConn database.DB

type ServiceExample struct {
	Log golog.MyLogger
	//Store       Storage
	dbConn      database.DB
	JwtSecret   []byte
	JwtDuration int
}

// login is just a trivial stupid example to test this server
// you should use the jwt token returned from LoginUser  in github.com/lao-tseu-is-alive/go-cloud-k8s-user-group'
// and share the same secret with the above component
func (s ServiceExample) login(ctx echo.Context) error {
	s.Log.Debug("entering login() \n##request: %+v", ctx.Request())

	uLogin := new(users.UserLogin)
	username := ctx.FormValue("login")
	fakePassword := ctx.FormValue("pass")
	s.Log.Debug("username: %s , password: %s", username, fakePassword)
	// maybe it was not a form but a fetch data post
	if len(strings.Trim(username, " ")) < 1 {
		if err := ctx.Bind(uLogin); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user login or json format in request body: %v", err)
		}
	} else {
		uLogin.Username = username
		uLogin.PasswordHash = fakePassword
	}
	s.Log.Debug("About to check username: %s , password: %s", uLogin.Username, uLogin.PasswordHash)

	// JSON serialized userLogin
	jsonData, err := json.Marshal(uLogin)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error serializing user login: %v", err))
	}

	/*
	userObjStore, err := users.GetStorageInstance("pgx", dbConn, s.Log)
	if err != nil {
		s.Log.Fatal("💥💥 error doing users.GetStorageInstance error: %v'\n", err)
	}
	loginService := users.Service{
		Log: s.Log,
		Store: userObjStore,
		JwtSecret: []byte(s.JwtSecret),
		JwtDuration: s.JwtDuration,
	}
	//e.POST("/login", loginService.LoginUser)
	
 	userSvcResp := loginService.LoginUser(ctx)	
	s.Log.Debug("response: %+v", userSvcResp)
	
	return userSvcResp
	*/
	
	//HTTP call to user service
	urlSvc := os.Getenv("GO_USER_SVC_URL")
	req, err := http.NewRequest("POST", urlSvc + "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error creating request to user service: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")

	// Request with context
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error calling user service: %v", err))
	}
	defer resp.Body.Close()

	// Reading response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error reading response from user service: %v", err))
	} else {
		// Verify response status
		if resp.StatusCode != http.StatusOK {
			//bodyStr := string(body)
			var message struct {
				Message string `json:"message"`
			}
			err = json.Unmarshal(body, &message)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error unmarshalling response from user service: %v", err))
			}
			return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error checking response status: %v %v", resp.StatusCode, message.Message))
		}
	}

	// Unmarshal response
	var result interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error unmarshalling response from user service: %v", err))
	}

	// Use of JWT token
	var token string
	switch v := result.(type) {
	case map[string]interface{}:
		token = v["token"].(string)
	case string:
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error parsing response from user service: %s", v))
	}
	s.Log.Debug("JWT token: %s", token)

	/*
	// Throws unauthorized error
	if username != defaultUsername || fakePassword != defaultFakeStupidPass {
		return ctx.JSON(http.StatusUnauthorized, "username not found or password invalid")
	}
	*/

	s.Log.Info("LoginUser(%s) succesfull login", uLogin.Username)
	return ctx.JSON(http.StatusOK, echo.Map{
		"token": token,
	})


	/*
	// Set custom claims
	claims := &goserver.JwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "",
			Audience:  nil,
			Issuer:    "",
			Subject:   "",
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * time.Duration(s.JwtDuration))},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			NotBefore: nil,
		},
		Id:       999,
		Name:     "Bill Whatever",
		Email:    "bill@whatever.com",
		Username: defaultUsername,
		IsAdmin:  false,
	}

	// Create token with claims
	signer, _ := jwt.NewSignerHS(jwt.HS512, s.JwtSecret)
	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		return err
	}
	s.Log.Info("LoginUser(%s) succesfull login for user id (%d)", claims.Username, claims.Id)
	return ctx.JSON(http.StatusOK, echo.Map{
		"token": token.String(),
	})
	*/

}

func (s ServiceExample) restricted(ctx echo.Context) error {
	s.Log.Debug("trace: entering restricted zone()")
	// get the current user from JWT TOKEN
	u := ctx.Get("jwtdata").(*jwt.Token)
	claims := goserver.JwtCustomClaims{}
	err := u.DecodeClaims(&claims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	//callerUserId := claims.Id
	// you can check if the user is not active anymore and RETURN 401 Unauthorized
	//if !s.Store.IsUserActive(currentUserId) {
	//	return echo.NewHTTPError(http.StatusUnauthorized, "current calling user is not active anymore")
	//}
	return ctx.JSON(http.StatusCreated, claims)
}

func isDBAlive() bool {
	dbVer, err := dbConn.GetVersion()
	if err != nil {
		return false
	}
	if len(dbVer) < 2 {
		return false
	}
	return true
}

func checkReady(string) bool {
	// we decide what makes us ready, is a valid  connection to the database
	if !isDBAlive() {
		return false
	}
	return true
}

func checkHealthy(string) bool {
	// you decide what makes you ready, may be it is the connection to the database
	//if !isDBAlive() {
	//	return false
	//}
	return true
}

func main() {
	prefix := fmt.Sprintf("%s ", version.APP)
	l, err := golog.NewLogger("zap", golog.DebugLevel, prefix)
	if err != nil {
		log.Fatalf("💥💥 error log.NewLogger error: %v\n", err)
	}
	l.Debug("Starting %s v:%s", version.APP, version.VERSION)
	l.Info("Starting %s v:%s", version.APP, version.VERSION)
	l.Warn("Starting %s v:%s", version.APP, version.VERSION)
	l.Error("Starting %s v:%s", version.APP, version.VERSION)
	secret, err := config.GetJwtSecretFromEnv()
	if err != nil {
		l.Fatal("💥💥 error doing config.GetJwtSecretFromEnv() error: %v\n", err)
	}
	tokenDuration, err := config.GetJwtDurationFromEnv(60)
	if err != nil {
		l.Fatal("💥💥 error doing config.GetJwtDurationFromEnv() error: %v\n", err)
	}
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
		tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	if err != nil {
		l.Fatal("💥💥 error doing config.GetPgDbDsnUrlFromEnv error: %v\n", err)
	}
	dbConn, err = database.GetInstance("pgx", dbDsn, runtime.NumCPU(), l)
	if err != nil {
		l.Fatal("💥💥 error doing database.GetInstance(\"pgx\", dbDsn)  : %v\n", err)
	}
	defer dbConn.Close()

	// example of go-migrate db migration with embed files in go program
	// https://github.com/golang-migrate/migrate
	// https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
	d, err := iofs.New(sqlMigrations, defaultSqlDbMigrationsPath)
	if err != nil {
		l.Fatal("💥💥 error doing iofs.New for db migrations  error: %v\n", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, strings.Replace(dbDsn, "postgres", "pgx", 1))
	if err != nil {
		l.Fatal("💥💥 error doing migrate.NewWithSourceInstance(iofs, dbURL:%s)  error: %v\n", dbDsn, err)
	}
	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			l.Fatal("💥💥 error doing migrate.Up error: %v\n", err)
		}
	}

	yourService := ServiceExample{
		Log:         l,
		dbConn:      dbConn,
		JwtSecret:   []byte(secret),
		JwtDuration: tokenDuration,
	}

	listenAddr, err := config.GetPortFromEnv(defaultPort)
	if err != nil {
		l.Fatal("💥💥 error calling GetPortFromEnv error: %v'\n", err)
	}
	l.Info("'Will start HTTP server listening on port %s'", listenAddr)
	server := goserver.NewGoHttpServer(listenAddr, l, defaultWebRootDir, content, defaultSecuredApi)
	e := server.GetEcho()
	e.GET("/readiness", server.GetReadinessHandler(checkReady, "Connection to DB"))
	e.GET("/health", server.GetHealthHandler(checkHealthy, "Connection to DB"))
	// Login route
	e.POST("/login", yourService.login)
	r := server.GetRestrictedGroup()
	r.GET("/secret", yourService.restricted)

	objStore, err := trees.GetStorageInstance("pgx", dbConn, l)
	if err != nil {
		l.Fatal("💥💥 error doing trees.GetStorageInstance error: %v'\n", err)
	}
	// now with restricted group reference you can register your secured handlers defined in OpenApi objects.yaml
	objService := trees.Service{
		Log:         l,
		Store:       objStore,
		JwtSecret:   []byte(secret),
		JwtDuration: tokenDuration,
	}
	trees.RegisterHandlers(r, &objService)

	loginExample := fmt.Sprintf("curl -v -X POST -d 'login=%s' -d 'pass=%s' http://localhost%s/login", defaultUsername, defaultFakeStupidPass, listenAddr)
	getSecretExample := fmt.Sprintf(" curl -v  -H \"Authorization: Bearer ${TOKEN}\" http://localhost%s%s/secret |jq\n", listenAddr, defaultSecuredApi)
	l.Info("From another terminal just try :\n %s", loginExample)
	l.Info("Then type export TOKEN=your_token_above_goes_here   \n %s", getSecretExample)

	err = server.StartServer()
	if err != nil {
		l.Error("💥💥 ERROR: 'calling echo.Start(%s) got error: %v'\n", listenAddr, err)
	}

}
