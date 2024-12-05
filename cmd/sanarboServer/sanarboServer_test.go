package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
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
		"password_hash": "58a38e08c4e6361b1f79c6c64417e7a2aed575a5158ec7d394e5d6f5f3a0d0c9", // Utiliser le mot de passe correspondant au hash stocké
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
