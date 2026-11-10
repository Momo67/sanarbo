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
	// Get, GetBuildingCenter and GetGestionComEmplacementsCentroidEmplacementId return getErr when set
	getErr error
	// SearchTreesByName returns searchErr and records its arguments
	searchErr   error
	searchCalls int
	searchText  string
	searchLimit int
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
func (f *fakeStorage) GetBuildingCenter(id int32) (*Center, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &Center{Geometry: "POINT(2538000 1152000)"}, nil
}
func (f *fakeStorage) GetGestionComEmplacementsCentroidEmplacementId(id int32) (*EmplacementCentroid, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &EmplacementCentroid{Geometry: "POINT(2538000 1152000)", Surface: 120}, nil
}
func (f *fakeStorage) SearchTreesByName(text string, limit int) ([]*TreeList, error) {
	f.searchCalls++
	f.searchText = text
	f.searchLimit = limit
	if f.searchErr != nil {
		return nil, f.searchErr
	}
	return []*TreeList{}, nil
}
func (f *fakeStorage) Get(id int32) (*Tree, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &Tree{Id: id, Name: "Platane 42"}, nil
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

func TestServiceGet(t *testing.T) {
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_trees")
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}

	tests := []struct {
		name           string
		getErr         error
		wantStatusCode int
	}{
		{name: "existing tree returns 200", wantStatusCode: http.StatusOK},
		{name: "unknown tree returns 404", getErr: ErrNoRecordFound, wantStatusCode: http.StatusNotFound},
		{name: "database error returns 500", getErr: errors.New("connection lost"), wantStatusCode: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				Log:   l,
				Store: &fakeStorage{getErr: tt.getErr},
			}
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/goapi/v1/trees/42", nil)
			rec := httptest.NewRecorder()

			err := s.Get(e.NewContext(req, rec), 42)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatusCode, rec.Code)
		})
	}
}

func TestServiceCenterAndCentroidNotFound(t *testing.T) {
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_trees")
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}
	handlers := map[string]func(s Service, ctx echo.Context) error{
		"GetBuildingCenter": func(s Service, ctx echo.Context) error { return s.GetBuildingCenter(ctx, 42) },
		"GetGestionComEmplacementsCentroidEmplacementId": func(s Service, ctx echo.Context) error {
			return s.GetGestionComEmplacementsCentroidEmplacementId(ctx, 42)
		},
	}
	tests := []struct {
		name           string
		getErr         error
		wantStatusCode int
	}{
		{name: "found returns 200", wantStatusCode: http.StatusOK},
		{name: "not found returns 404", getErr: ErrNoRecordFound, wantStatusCode: http.StatusNotFound},
		{name: "database error returns 500", getErr: errors.New("connection lost"), wantStatusCode: http.StatusInternalServerError},
	}
	for handlerName, handler := range handlers {
		for _, tt := range tests {
			t.Run(handlerName+" "+tt.name, func(t *testing.T) {
				s := Service{Log: l, Store: &fakeStorage{getErr: tt.getErr}}
				e := echo.New()
				rec := httptest.NewRecorder()

				err := handler(s, e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), rec))

				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatusCode, rec.Code)
			})
		}
	}
}

func TestServiceSearchTreesByName(t *testing.T) {
	l, err := golog.NewLogger("simple", golog.ErrorLevel, "test_trees")
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}

	tests := []struct {
		name           string
		query          string
		searchErr      error
		wantStatusCode int
		wantText       string
		wantLimit      int
	}{
		{name: "default limit is 100", query: "name=Platane", wantStatusCode: http.StatusOK, wantText: "Platane", wantLimit: 100},
		{name: "text is trimmed and limit is passed", query: "name=%20%20Culotte%20suisse%20&limit=5", wantStatusCode: http.StatusOK, wantText: "Culotte suisse", wantLimit: 5},
		{name: "two accented characters are enough", query: "name=%C3%A9%C3%A9", wantStatusCode: http.StatusOK, wantText: "éé", wantLimit: 100},
		{name: "missing name is rejected", query: "", wantStatusCode: http.StatusBadRequest},
		{name: "one character is rejected", query: "name=P", wantStatusCode: http.StatusBadRequest},
		{name: "blank name is rejected", query: "name=%20%20%20", wantStatusCode: http.StatusBadRequest},
		{name: "limit 0 is rejected", query: "name=Platane&limit=0", wantStatusCode: http.StatusBadRequest},
		{name: "limit above 1000 is rejected", query: "name=Platane&limit=1001", wantStatusCode: http.StatusBadRequest},
		{name: "database error returns 500", query: "name=Platane", searchErr: errors.New("connection lost"), wantStatusCode: http.StatusInternalServerError, wantText: "Platane", wantLimit: 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &fakeStorage{searchErr: tt.searchErr}
			// go through the generated router, to check /trees/search is not routed to /trees/{treeId}
			e := echo.New()
			RegisterHandlers(e, &Service{Log: l, Store: store})
			req := httptest.NewRequest(http.MethodGet, "/trees/search?"+tt.query, nil)
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatusCode, rec.Code, rec.Body.String())
			if tt.wantText == "" {
				assert.Equal(t, 0, store.searchCalls, "Store.SearchTreesByName must not be called")
				return
			}
			assert.Equal(t, 1, store.searchCalls)
			assert.Equal(t, tt.wantText, store.searchText)
			assert.Equal(t, tt.wantLimit, store.searchLimit)
			if tt.wantStatusCode == http.StatusOK {
				assert.JSONEq(t, "[]", rec.Body.String())
			}
		})
	}
}
