package trees

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goHttpEcho"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/stretchr/testify/assert"
)

// fakeJwtChecker returns fixed claims without parsing any token.
// Only GetJwtCustomClaimsFromContext is implemented, other methods will panic if called.
type fakeJwtChecker struct {
	goHttpEcho.JwtChecker
	user *goHttpEcho.UserInfo
}

func (f fakeJwtChecker) GetJwtCustomClaimsFromContext(c echo.Context) *goHttpEcho.JwtCustomClaims {
	return &goHttpEcho.JwtCustomClaims{User: f.user}
}

// fakeStorage implements only the Storage methods used by Service.Update, other methods will panic if called.
type fakeStorage struct {
	Storage
	isAdmin     bool
	isEditor    bool
	updateCalls int
}

func (f *fakeStorage) IsObjectAdmin(id int32) bool  { return f.isAdmin }
func (f *fakeStorage) IsObjectEditor(id int32) bool { return f.isEditor }
func (f *fakeStorage) Exist(id int32) bool          { return true }
func (f *fakeStorage) Update(id int32, object Tree) (*Tree, error) {
	f.updateCalls++
	object.Id = id
	return &object, nil
}

func TestServiceUpdatePermissions(t *testing.T) {
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_trees")
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}
	validTree := `{"id": 42, "name": "Platane 42", "geom": "POINT(2538000 1152000)", "tree_attributes": {"circonference": 120}}`

	tests := []struct {
		name            string
		isAdmin         bool
		isEditor        bool
		wantStatusCode  int
		wantUpdateCalls int
	}{
		{name: "user without admin or editor privilege is forbidden", wantStatusCode: http.StatusForbidden, wantUpdateCalls: 0},
		{name: "editor can update", isEditor: true, wantStatusCode: http.StatusOK, wantUpdateCalls: 1},
		{name: "admin can update", isAdmin: true, wantStatusCode: http.StatusOK, wantUpdateCalls: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStorage{isAdmin: tt.isAdmin, isEditor: tt.isEditor}
			s := Service{
				Log:   l,
				Store: store,
				Server: &goHttpEcho.Server{
					JwtCheck: fakeJwtChecker{user: &goHttpEcho.UserInfo{UserId: 7, ExternalId: 1007}},
				},
			}
			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/goapi/v1/trees/42", strings.NewReader(validTree))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			err := s.Update(e.NewContext(req, rec), 42)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
			assert.Equal(t, tt.wantUpdateCalls, store.updateCalls, "Store.Update calls")
		})
	}
}
