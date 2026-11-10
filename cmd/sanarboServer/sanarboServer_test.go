package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goHttpEcho"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/stretchr/testify/assert"
)

const (
	DEBUG                           = true
	assertCorrectStatusCodeExpected = "expected status code should be returned"
)

type testStruct struct {
	name           string
	contentType    string
	wantStatusCode int
	wantBody       string
	paramKeyValues map[string]string
	httpMethod     string
	url            string
	body           string
}

func TestServiceRestricted(t *testing.T) {
	type fields struct {
		Log         golog.MyLogger
		dbConn      database.DB
		JwtSecret   []byte
		JwtDuration int
	}
	type args struct {
		ctx echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				Logger: tt.fields.Log,
				dbConn: tt.fields.dbConn,
				server: &goHttpEcho.Server{
					Authenticator: nil,
					JwtCheck:      nil,
					VersionReader: nil,
				},
			}
			if err := s.restricted(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("restricted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckHealthy(t *testing.T) {
	type args struct {
		info string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkHealthy(tt.args.info); got != tt.want {
				t.Errorf("checkHealthy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func waitForServer(addr string, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("server did not start within %s", timeout)
		}
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		conn.Close()
		return nil
	}
}

// TestMainExec is instantiating the "real" main code using the env variable (in your .env.development.local files if you use the Makefile rule)
func TestMainExec(t *testing.T) {
	listenPort := config.GetPortFromEnvOrPanic(defaultPort)
	listenIP := config.GetListenIpFromEnvOrPanic("0.0.0.0")
	listenAddr := fmt.Sprintf("%s://%s:%d", "http", listenIP, listenPort)
	fmt.Printf("INFO: 'Will start HTTP server listening on port %s'\n", listenAddr)

	newRequest := func(method, url string, body string, contenttype string) *http.Request {
		fmt.Printf("INFO: 💥💥'newRequest %s on %s ##BODY : %+v'\n", method, url, body)
		r, err := http.NewRequest(method, url, strings.NewReader(body))
		if err != nil {
			t.Fatalf("### ERROR http.NewRequest %s on [%s] error is :%v\n", method, url, err)
		}
		/*
		if method == http.MethodPost {
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		*/
		r.Header.Set("Content-Type", contenttype)
		return r
	}

	// Créer l'objet JSON pour le login
	loginData := map[string]string{
		"username":      "go-admin",
		"password_hash": "db48f76f6c362ce9c767f53d9e0dfc10f67690e9f72d23cab0454a90a56acc76", // Utiliser le mot de passe correspondant au hash stocké
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	tests := []testStruct{
		{
			name:           "1: Get on default get handler should contain html tag",
			wantStatusCode: http.StatusOK,
			contentType:    MIMEHtmlCharsetUTF8,
			wantBody:       "<html",
			paramKeyValues: make(map[string]string, 0),
			httpMethod:     http.MethodGet,
			url:            "/",
			body:           "",
		},
		{
			name:           "2: Post on default get handler should return an http error method not allowed ",
			wantStatusCode: http.StatusMethodNotAllowed,
			contentType:    MIMEHtmlCharsetUTF8,
			wantBody:       "Method Not Allowed",
			paramKeyValues: make(map[string]string, 0),
			httpMethod:     http.MethodPost,
			url:            "/",
			body:           `{"junk":"test with junk text"}`,
		},
		{
			name:           "3: Get on nonexistent route should return an http error not found ",
			wantStatusCode: http.StatusNotFound,
			contentType:    MIMEHtmlCharsetUTF8,
			wantBody:       "page not found",
			paramKeyValues: make(map[string]string, 0),
			httpMethod:     http.MethodGet,
			url:            "/aroutethatwillneverexisthere",
			body:           "",
		},
		{
			name:           "4: POST to login with valid credential should return a JWT token ",
			wantStatusCode: http.StatusOK,
			contentType:    echo.MIMEApplicationJSON,
			wantBody:       "token",
			paramKeyValues: make(map[string]string, 0),
			httpMethod:     http.MethodPost,
			url:            "/login",
			body:           string(jsonData),
		},
	}

	// starting main in his own go routine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		main()
	}()

	err = waitForServer(fmt.Sprintf("0.0.0.0:%d", listenPort), 2*time.Second)
	if err != nil {
		t.Fatalf("Server did not start in time: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := newRequest(tt.httpMethod, listenAddr+tt.url, tt.body, tt.contentType)
			//r.Header.Set(HeaderContentType, tt.contentType)
			resp, err := http.DefaultClient.Do(r)
			if DEBUG {
				fmt.Printf("### %s : %s on %s\n", tt.name, r.Method, r.URL)
			}
			if err != nil {
				fmt.Printf("### GOT ERROR : %s\n%+v", err, resp)
				t.Fatal(err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode, assertCorrectStatusCodeExpected)
			receivedJson, _ := io.ReadAll(resp.Body)

			if DEBUG {
				fmt.Printf("WANTED   :%T - %#v\n", tt.wantBody, tt.wantBody)
				fmt.Printf("RECEIVED :%T - %#v\n", receivedJson, string(receivedJson))
			}
			// check that receivedJson contains the specified tt.wantBody substring . https://pkg.go.dev/github.com/stretchr/testify/assert#Contains
			assert.Contains(t, string(receivedJson), tt.wantBody, "Response should contain what was expected.")
		})
	}
}

func TestRedactDsnPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		dsn      string
		want     string
	}{
		{
			name:     "simple password",
			password: "secret",
			dsn:      "postgres://sanarbo:secret@127.0.0.1:5432/sanarbo?sslmode=prefer",
			want:     "postgres://sanarbo:xxxxx@127.0.0.1:5432/sanarbo?sslmode=prefer",
		},
		{
			name:     "password with url special characters is not url-escaped by config",
			password: "p/a?s#s@w:rd",
			dsn:      "postgres://sanarbo:p/a?s#s@w:rd@db.local:5432/sanarbo?sslmode=require",
			want:     "postgres://sanarbo:xxxxx@db.local:5432/sanarbo?sslmode=require",
		},
		{
			name: "dsn without password is unchanged",
			dsn:  "postgres://sanarbo@127.0.0.1:5432/sanarbo",
			want: "postgres://sanarbo@127.0.0.1:5432/sanarbo",
		},
		{
			name:     "unexpected format is not shown at all",
			password: "secret",
			dsn:      "host=127.0.0.1 user=sanarbo password=secret",
			want:     "[dsn not shown]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := redactDsnPassword(tt.dsn)
			assert.Equal(t, tt.want, got)
			if tt.password != "" {
				assert.NotContains(t, got, tt.password, "password must not appear in redacted dsn")
			}
		})
	}
}

func TestGetUrlFromEnvOrPanic(t *testing.T) {
	const envName = "SANARBO_TEST_URL"
	tests := []struct {
		name      string
		envValue  *string
		def       string
		want      string
		wantPanic bool
	}{
		{name: "unset env uses relative path default", envValue: nil, def: "/login", want: "/login"},
		{name: "unset env uses root path default", envValue: nil, def: "/", want: "/"},
		{name: "empty env uses default (Dockerfile defines empty variables)", envValue: strPtr(""), def: "/login", want: "/login"},
		{name: "blank env uses default", envValue: strPtr("  "), def: "/", want: "/"},
		{name: "absolute http url", envValue: strPtr("http://ingress.local/user-group"), def: "/login", want: "http://ingress.local/user-group"},
		{name: "absolute https url", envValue: strPtr("https://sanarbo.lausanne.ch"), def: "/", want: "https://sanarbo.lausanne.ch"},
		{name: "path", envValue: strPtr("/user-group/login"), def: "/", want: "/user-group/login"},
		{name: "url without scheme", envValue: strPtr("ingress.local/user-group"), def: "/", wantPanic: true},
		{name: "host and port without scheme", envValue: strPtr("localhost:9999"), def: "/", wantPanic: true},
		{name: "unsupported scheme", envValue: strPtr("ftp://ingress.local"), def: "/", wantPanic: true},
		{name: "scheme without host", envValue: strPtr("http://"), def: "/", wantPanic: true},
		{name: "protocol-relative url", envValue: strPtr("//ingress.local"), def: "/", wantPanic: true},
		{name: "invalid default", envValue: nil, def: "login", wantPanic: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Setenv restores the previous value at the end of the test
			t.Setenv(envName, "")
			if tt.envValue == nil {
				os.Unsetenv(envName)
			} else {
				os.Setenv(envName, *tt.envValue)
			}
			if tt.wantPanic {
				assert.Panics(t, func() { getUrlFromEnvOrPanic(envName, tt.def) })
				return
			}
			assert.Equal(t, tt.want, getUrlFromEnvOrPanic(envName, tt.def))
		})
	}
}

// TestGetInfoUrlsDefaults checks the defaults used by GetInfoHandlerOrPanic no longer panic when the env variables are not set
func TestGetInfoUrlsDefaults(t *testing.T) {
	t.Setenv("GO_USER_SVC_URL", "")
	t.Setenv("GO_BACKEND_URL", "")
	os.Unsetenv("GO_USER_SVC_URL")
	os.Unsetenv("GO_BACKEND_URL")

	assert.Equal(t, "/login", GetAuthenticationUrlFromEnvOrPanic("/login"))
	assert.Equal(t, "/", GetBackendUrlFromEnvOrPanic("/"))
	assert.NotPanics(t, func() { GetInfoHandlerOrPanic() })
}

func strPtr(s string) *string { return &s }

func TestGetLogLevelFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envValue *string
		want     golog.Level
		wantErr  bool
	}{
		{name: "unset uses default info", envValue: nil, want: golog.InfoLevel},
		{name: "empty uses default info", envValue: strPtr(""), want: golog.InfoLevel},
		{name: "trace", envValue: strPtr("trace"), want: golog.TraceLevel},
		{name: "debug", envValue: strPtr("debug"), want: golog.DebugLevel},
		{name: "info", envValue: strPtr("info"), want: golog.InfoLevel},
		{name: "warn", envValue: strPtr("warn"), want: golog.WarnLevel},
		{name: "error", envValue: strPtr("error"), want: golog.ErrorLevel},
		{name: "fatal", envValue: strPtr("fatal"), want: golog.FatalLevel},
		{name: "case and spaces are ignored", envValue: strPtr("  WARN "), want: golog.WarnLevel},
		{name: "unknown level is an error", envValue: strPtr("verbose"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Setenv restores the previous value at the end of the test
			t.Setenv("LOG_LEVEL", "")
			if tt.envValue == nil {
				os.Unsetenv("LOG_LEVEL")
			} else {
				os.Setenv("LOG_LEVEL", *tt.envValue)
			}
			got, err := getLogLevelFromEnv(defaultLogLevel)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// fakeJwtChecker returns fixed claims without parsing any token, other methods will panic if called.
type fakeJwtChecker struct {
	goHttpEcho.JwtChecker
	user *goHttpEcho.UserInfo
}

func (f fakeJwtChecker) GetJwtCustomClaimsFromContext(c echo.Context) *goHttpEcho.JwtCustomClaims {
	return &goHttpEcho.JwtCustomClaims{User: f.user}
}

func TestRestrictedReturns200(t *testing.T) {
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_main")
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}
	s := Service{
		Logger: l,
		server: &goHttpEcho.Server{JwtCheck: fakeJwtChecker{user: &goHttpEcho.UserInfo{UserId: 7, Login: "test"}}},
	}
	e := echo.New()
	rec := httptest.NewRecorder()

	err = s.restricted(e.NewContext(httptest.NewRequest(http.MethodGet, "/goapi/v1/secret", nil), rec))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code, "a GET returning the claims creates nothing, it must not answer 201")
	assert.Contains(t, rec.Body.String(), `"login":"test"`)
}
