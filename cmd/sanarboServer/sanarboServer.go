package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goHttpEcho"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/metadata"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"github.com/lao-tseu-is-alive/sanarbo/pkg/trees"
	"github.com/lao-tseu-is-alive/sanarbo/pkg/version"
)

const (
	APP                        = "sanarbo"
	defaultPort                = 9999
	defaultDBPort              = 5432
	defaultDBIp                = "127.0.0.1"
	defaultDBSslMode           = "prefer"
	defaultReadTimeout         = 10 * time.Second // max time to read request from the client
	defaultWebRootDir          = "sanarboFront/dist/"
	defaultSqlDbMigrationsPath = "db/migrations"
	defaultSecuredApi          = "/goapi/v1"
	defaultAdminUser           = "goadmin"
	defaultAdminEmail          = "goadmin@yourdomain.org"
	defaultAdminId             = 960901
	charsetUTF8                = "charset=UTF-8"
	MIMEHtml                   = "text/html"
	MIMEHtmlCharsetUTF8        = MIMEHtml + "; " + charsetUTF8
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

type Service struct {
	Logger golog.MyLogger
	dbConn database.DB
	server *goHttpEcho.Server
}

/*
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

	//HTTP call to user service
	urlSvc := os.Getenv("GO_USER_SVC_URL")
	urlSvc = strings.Replace(urlSvc, "\"", "", -1)
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
	var tokenString string
	switch v := result.(type) {
	case map[string]interface{}:
		tokenString = v["token"].(string)
	case string:
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error parsing response from user service: %s", v))
	}
	s.Log.Debug("JWT token: %s", tokenString)

	verifier, err := jwt.NewVerifierHS(jwt.HS512, []byte(s.JwtSecret))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error creating verifier: %v", err))
	}

	token, err := jwt.Parse([]byte(tokenString), verifier)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error parsing token: %v", err))
	}

	claims := goHttpEcho.JwtCustomClaims{}
	err = token.DecodeClaims(&claims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error decoding claims: %v", err))
	}
	s.Log.Debug("Id: %d, Name: %s, Email: %s, Username: %s, IsAdmin: %t", claims.Id, claims.Name, claims.Email, claims.Username, claims.IsAdmin)


	s.Log.Info("LoginUser(%s) succesfull login", uLogin.Username)
	return ctx.JSON(http.StatusOK, echo.Map{
		"token": tokenString,
	})
}
*/

func (s Service) restricted(ctx echo.Context) error {
	s.Logger.TraceHttpRequest("restricted", ctx.Request())
	// get the current user from JWT TOKEN
	claims := s.server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := claims.User.UserId
	s.Logger.Info("in restricted : currentUserId: %d", currentUserId)
	// you can check if the user is not active anymore and RETURN 401 Unauthorized
	//if !s.Store.IsUserActive(currentUserId) {
	//	return echo.NewHTTPError(http.StatusUnauthorized, "current calling user is not active anymore")
	//}
	return ctx.JSON(http.StatusCreated, claims)
}

func checkHealthy(info string) bool {
	// you decide what makes you ready, may be it is the connection to the database
	//if !stillConnectedToDB {
	//	return false
	//}
	return true
}

// GetAuthenticationUrlFromEnvOrPanic returns the authentication url to be used for JWT the content of the env variable:
// GO_USER_SVC_URL : string containing the authentication url to use for JWT authentication
func GetAuthenticationUrlFromEnvOrPanic(defaultAuthenticationUrl string) string {
	authenticationUrl := defaultAuthenticationUrl
	val, exist := os.LookupEnv("GO_USER_SVC_URL")
	if exist {
		authenticationUrl = val
	}
	if utf8.RuneCountInString(authenticationUrl) < 10 {
		panic(fmt.Sprintf("💥💥 ERROR: CONFIG ENV GO_USER_SVC_URL should contain at least 10 characters (got %d).",
			utf8.RuneCountInString(val)))
	}
	return fmt.Sprintf("%s", authenticationUrl)
}

// GetBackendUrlFromEnvOrPanic returns the backend url to be used for JWT the content of the env variable:
// GO_BACKEND_URL : string containing the backend url to use for API calls
func GetBackendUrlFromEnvOrPanic(defaultBackendUrl string) string {
	backendUrl := defaultBackendUrl
	val, exist := os.LookupEnv("GO_BACKEND_URL")
	if exist {
		backendUrl = val
	}
	if utf8.RuneCountInString(backendUrl) < 10 {
		panic(fmt.Sprintf("💥💥 ERROR: CONFIG ENV GO_BACKEND_URL should contain at least 10 characters (got %d).",
			utf8.RuneCountInString(val)))
	}
	return fmt.Sprintf("%s", backendUrl)
}

func GetInfoHandlerOrPanic() echo.HandlerFunc {
	myBackendUrl := GetBackendUrlFromEnvOrPanic("/")
	myAuthenticationUrl := GetAuthenticationUrlFromEnvOrPanic("/login")
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"app":                APP,
			"version":            version.VERSION,
			"repo":               version.REPOSITORY,
			"backend_url":        myBackendUrl,
			"authentication_url": myAuthenticationUrl,
			"restricted_url":     defaultSecuredApi,
		})
	}
}


func main() {
	l, err := golog.NewLogger("zap", golog.DebugLevel, APP)
	if err != nil {
		log.Fatalf("💥💥 error log.NewLogger error: %v'\n", err)
	}
	l.Info("🚀🚀 Starting App:'%s', ver:%s, from: %s", APP, version.VERSION, version.REPOSITORY)

	dbDsn := config.GetPgDbDsnUrlFromEnvOrPanic(defaultDBIp, defaultDBPort, tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	db, err := database.GetInstance("pgx", dbDsn, runtime.NumCPU(), l)
	if err != nil {
		l.Fatal("💥💥 error doing database.GetInstance(pgx ...) error: %v", err)
	}
	defer db.Close()

	dbVersion, err := db.GetVersion()
	if err != nil {
		l.Fatal("💥💥 error doing dbConn.GetVersion() error: %v", err)
	}
	l.Info("connected to db version : %s", dbVersion)

	// checking metadata information
	metadataService := metadata.Service{Log: l, Db: db}
	metadataService.CreateMetadataTableOrFail()
	found, ver := metadataService.GetServiceVersionOrFail(version.APP)
	if found {
		l.Info("service %s was found in metadata with version: %s", version.APP, ver)
	} else {
		l.Info("service %s was not found in metadata", version.APP)
	}
	metadataService.SetServiceVersionOrFail(version.APP, version.VERSION)

	// https://github.com/golang-migrate/migrate
	d, err := iofs.New(sqlMigrations, defaultSqlDbMigrationsPath)
	if err != nil {
		l.Fatal("💥💥 error doing iofs.New for db migrations  error: %v\n", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, strings.Replace(dbDsn, "postgres", "pgx5", 1))
	if err != nil {
		l.Fatal("💥💥 error doing migrate.NewWithSourceInstance(iofs, dbURL:%s)  error: %v\n", dbDsn, err)
	}

	err = m.Up()
	if err != nil {
		//if err == m.
		if !errors.Is(err, migrate.ErrNoChange) {
			l.Fatal("💥💥 error doing migrate.Up error: %v\n", err)
		}
	}

	myVersionReader := goHttpEcho.NewSimpleVersionReader(APP, version.VERSION, version.REPOSITORY)
	// Create a new JWT checker
	myJwt := goHttpEcho.NewJwtChecker(
		config.GetJwtSecretFromEnvOrPanic(),
		config.GetJwtIssuerFromEnvOrPanic(),
		APP,
		config.GetJwtContextKeyFromEnvOrPanic(),
		config.GetJwtDurationFromEnvOrPanic(60),
		l)
	// Create a new Authenticator with a simple admin user
	myAuthenticator := goHttpEcho.NewSimpleAdminAuthenticator(&goHttpEcho.UserInfo{
		UserId:     config.GetAdminIdFromEnvOrPanic(defaultAdminId),
		ExternalId: config.GetAdminExternalIdFromEnvOrPanic(9999999),
		Name:       "NewSimpleAdminAuthenticator_Admin",
		Email:      config.GetAdminEmailFromEnvOrPanic(defaultAdminEmail),
		Login:      config.GetAdminUserFromEnvOrPanic(defaultAdminUser),
		IsAdmin:    false,
	},

		config.GetAdminPasswordFromEnvOrPanic(),
		myJwt)

	server := goHttpEcho.CreateNewServerFromEnvOrFail(
		defaultPort,
		"0.0.0.0", // defaultServerIp,
		&goHttpEcho.Config{
			ListenAddress: "",
			Authenticator: myAuthenticator,
			JwtCheck:      myJwt,
			VersionReader: myVersionReader,
			Logger:        l,
			WebRootDir:    defaultWebRootDir,
			Content:       content,
			RestrictedUrl: "/goapi/v1",
		},
	)

	userStore := users.GetStorageInstanceOrPanic("pgx", db, l)

	userService := users.Service{
		Logger: l,
		DbConn: db,
		Store:  userStore,
		Server: server,
	}

	e := server.GetEcho()
	e.POST("/login", userService.LoginUser)
	e.GET("/readiness", server.GetReadinessHandler(func(info string) bool {
		ver, err := db.GetVersion()
		if err != nil {
			l.Error("Error getting db version : %v", err)
			return false
		}
		l.Info("Connected to DB version : %s", ver)
		return true
	}, "Connection to DB"))
	e.GET("/health", server.GetHealthHandler(checkHealthy, "Connection to DB"))
	e.GET("/info", GetInfoHandlerOrPanic())

	yourService := Service{
		Logger: l,
		dbConn: db,
		server: server,
	}

	r := server.GetRestrictedGroup()
	users.RegisterHandlers(r, &userService)

	r.GET("/secret", yourService.restricted)
	r.GET("/status", userService.GetStatus)
	r.GET("/users/maxid", userService.GetMaxId)

	objStore, err := trees.GetStorageInstance("pgx", db, l)
	if err != nil {
		l.Fatal("💥💥 error doing trees.GetStorageInstance error: %v'\n", err)
	}
	// now with restricted group reference you can register your secured handlers defined in OpenApi objects.yaml
	objService := trees.Service{
		Log:    l,
		Store:  objStore,
		Server: server,
	}
	trees.RegisterHandlers(r, &objService)

	loginExample := fmt.Sprintf("curl -v -X POST -H \"Content-Type: application/json\" -d '{ \"username\": \"%s\", \"password_hash\": \"%s\" }' http://localhost:%d/login", "your_user", "your_password_hash", config.GetPortFromEnvOrPanic(defaultPort))
	getSecretExample := fmt.Sprintf(" curl -v  -H \"Authorization: Bearer ${TOKEN}\" http://localhost:%d/%s/secret |jq\n", config.GetPortFromEnvOrPanic(defaultPort), defaultSecuredApi)
	l.Info("From another terminal just try :\n %s", loginExample)
	l.Info("Then type export TOKEN=your_token_above_goes_here   \n %s", getSecretExample)

	err = server.StartServer()
	if err != nil {
		l.Fatal("💥💥 error doing server.StartServer error: %v'\n", err)
	}
}
