package trees

import (
	"errors"
	"fmt"
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

// fakeStorage implements only the Storage methods used by the tested handlers, other methods will panic if called.
type fakeStorage struct {
	Storage
	isAdmin     bool
	isEditor    bool
	updateCalls int
	created     *Tree
	// ValidateTrees returns validateErr and records its arguments
	validateErr        error
	validateCalls      int
	validatedIds       []int32
	validatedValidator int32
}

func (f *fakeStorage) IsObjectAdmin(id int32) bool  { return f.isAdmin }
func (f *fakeStorage) IsObjectEditor(id int32) bool { return f.isEditor }
func (f *fakeStorage) Exist(id int32) bool          { return true }
func (f *fakeStorage) Update(id int32, object Tree) (*Tree, error) {
	f.updateCalls++
	object.Id = id
	return &object, nil
}
func (f *fakeStorage) Create(object Tree) (*Tree, error) {
	f.created = &object
	return &object, nil
}
func (f *fakeStorage) IsObjectValidator(id int32) bool { return false }
func (f *fakeStorage) ValidateTrees(externalIds []int32, idValidator int32) error {
	f.validateCalls++
	f.validatedIds = externalIds
	f.validatedValidator = idValidator
	return f.validateErr
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

func TestServiceCreateIgnoresCreatorFromBody(t *testing.T) {
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_trees")
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}
	const authenticatedUserId = 7
	treeWithSpoofedCreator := `{"creator": 999, "name": "Platane 42", "geom": "POINT(2538000 1152000)", "tree_attributes": {"circonference": 120}}`

	store := &fakeStorage{isAdmin: true}
	s := Service{
		Log:   l,
		Store: store,
		Server: &goHttpEcho.Server{
			JwtCheck: fakeJwtChecker{user: &goHttpEcho.UserInfo{UserId: authenticatedUserId, ExternalId: 1007}},
		},
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/goapi/v1/trees", strings.NewReader(treeWithSpoofedCreator))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	err = s.Create(e.NewContext(req, rec))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	if assert.NotNil(t, store.created, "Store.Create should have been called") {
		assert.Equal(t, int32(authenticatedUserId), store.created.Creator, "creator must come from the JWT, not from the request body")
	}
}

func TestServiceSaveValidation(t *testing.T) {
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_trees")
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}
	const validatorExternalId = 1007

	tests := []struct {
		name              string
		body              string
		validateErr       error
		wantStatusCode    int
		wantValidateCalls int
		wantValidatedIds  []int32
	}{
		{
			name:              "only trees marked as validated are saved, in a single call",
			body:              `[{"external_id": 101, "is_validated": true}, {"external_id": 102, "is_validated": false}, {"external_id": 103, "is_validated": true}]`,
			wantStatusCode:    http.StatusOK,
			wantValidateCalls: 1,
			wantValidatedIds:  []int32{101, 103},
		},
		{
			name:              "an invalid entry anywhere in the list rejects the whole list before saving",
			body:              `[{"external_id": 101, "is_validated": true}, {"external_id": 0, "is_validated": true}]`,
			wantStatusCode:    http.StatusBadRequest,
			wantValidateCalls: 0,
		},
		{
			name:              "empty list is rejected",
			body:              `[]`,
			wantStatusCode:    http.StatusBadRequest,
			wantValidateCalls: 0,
		},
		{
			name:              "nothing marked as validated does not call the store",
			body:              `[{"external_id": 101, "is_validated": false}]`,
			wantStatusCode:    http.StatusOK,
			wantValidateCalls: 0,
		},
		{
			name:              "unknown tree returns 404",
			body:              `[{"external_id": 999999, "is_validated": true}]`,
			validateErr:       fmt.Errorf("%w : tree with external_id 999999", ErrNoRecordFound),
			wantStatusCode:    http.StatusNotFound,
			wantValidateCalls: 1,
			wantValidatedIds:  []int32{999999},
		},
		{
			name:              "database error returns 500 instead of a false success",
			body:              `[{"external_id": 101, "is_validated": true}]`,
			validateErr:       errors.New("connection lost"),
			wantStatusCode:    http.StatusInternalServerError,
			wantValidateCalls: 1,
			wantValidatedIds:  []int32{101},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStorage{isAdmin: true, validateErr: tt.validateErr}
			s := Service{
				Log:   l,
				Store: store,
				Server: &goHttpEcho.Server{
					JwtCheck: fakeJwtChecker{user: &goHttpEcho.UserInfo{UserId: 7, ExternalId: validatorExternalId}},
				},
			}
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/goapi/v1/validation", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			err := s.SaveValidation(e.NewContext(req, rec))

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
			assert.Equal(t, tt.wantValidateCalls, store.validateCalls, "Store.ValidateTrees calls")
			if tt.wantValidateCalls > 0 {
				assert.Equal(t, tt.wantValidatedIds, store.validatedIds)
				assert.Equal(t, int32(validatorExternalId), store.validatedValidator, "validator must come from the JWT")
			}
		})
	}
}
