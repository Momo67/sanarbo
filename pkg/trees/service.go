package trees

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goHttpEcho"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
)

const (
	noAdminPrivilege        = "current user has no admin privilege"
	noObjectAdminPrivilege  = "current user has no admin privilege on this object"
	noObjectEditorPrivilege = "current user has no editor privilege on this object"
	noObjectValidatorPrivilege = "current user has no validator privilege on this object"
)

type Service struct {
	Log    golog.MyLogger
	Store  Storage
	Server *goHttpEcho.Server
}


func (s Service) List(ctx echo.Context, params ListParams) error {
	s.Log.Debug("entering List() params:%v\n", params)

	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := claims.User.UserId
	s.Log.Info("in UserCreate : currentUserId: %d", currentUserId)

	var limit int = 100
	if params.Limit != nil {
		limit = int(*params.Limit)
	}
	var offset int = 0
	if params.Offset != nil {
		offset = int(*params.Offset)
	}
	list, err := s.Store.List(offset, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("there was a problem when calling store.List :%v", err))
	}
	return ctx.JSON(http.StatusOK, list)
}

func (s Service) Create(ctx echo.Context) error {
	s.Log.Debug("entering Create()")

	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := claims.User.UserId
	s.Log.Info("in UserCreate : currentUserId: %d", currentUserId)
	if !s.Store.IsObjectAdmin(int32(currentUserId)) {
		return ctx.JSON(http.StatusUnauthorized, noAdminPrivilege)
	}

	newTree := &Tree{
		Id:      0,
		Creator: int32(currentUserId),
	}
	if err := ctx.Bind(newTree); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Tree has invalid format [%v]", err))
	}
	if len(newTree.Name) < 1 {
		return ctx.JSON(http.StatusBadRequest, "Tree name cannot be empty")
	}
	if len(newTree.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Tree name minlength is 5 not (%d)", len(newTree.Name)))
	}
	if len(newTree.Geom) < 1 {
		return ctx.JSON(http.StatusBadRequest, "Tree geom cannot be empty")
	}
	if (TreeAttributes{}) == newTree.TreeAttributes {
		return ctx.JSON(http.StatusBadRequest, "Tree tree_attributes cannot be empty")
	}
	s.Log.Info("# Create() newTree : %#v\n", newTree)
	treeCreated, err := s.Store.Create(*newTree)
	if err != nil {
		msg := fmt.Sprintf("Create had an error saving tree:%#v, err:%#v", *newTree, err)
		s.Log.Error(msg)
		return ctx.JSON(http.StatusBadRequest, msg)
	}
	s.Log.Info("# Create() Tree %#v\n", treeCreated)
	return ctx.JSON(http.StatusCreated, treeCreated)
}

func (s Service) Delete(ctx echo.Context, objectId int32) error {
	s.Log.Debug("entering Delete(%d)\n", objectId)

	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := claims.User.UserId
	if !s.Store.IsObjectAdmin(int32(currentUserId)) {
		return ctx.JSON(http.StatusUnauthorized, noObjectAdminPrivilege)
	}
	if !s.Store.Exist(objectId) {
		msg := fmt.Sprintf("Delete(%d) cannot delete this id, it does not exist !", objectId)
		s.Log.Error(msg)
		return ctx.JSON(http.StatusNotFound, msg)
	} else {
		err := s.Store.Delete(objectId)
		if err != nil {
			msg := fmt.Sprintf("Delete(%d) got an error: %#v ", objectId, err)
			s.Log.Error(msg)
			return ctx.JSON(http.StatusInternalServerError, msg)
		}
		return ctx.NoContent(http.StatusNoContent)
	}
}

func (s Service) Get(ctx echo.Context, objectId int32) error {
	s.Log.Debug("entering Get(%d)", objectId)

	tree, err := s.Store.Get(objectId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving tree :%v", err))
	}
	return ctx.JSON(http.StatusOK, tree)
}

func (s Service) Update(ctx echo.Context, objectId int32) error {
	s.Log.Debug("entering Update(%d)\n", objectId)

	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	/*
		currentUserId := claims.User.UserId
		if !(s.Store.IsObjectAdmin(int32(currentUserId)) || s.Store.IsObjectEditor(int32(currentUserId))) {
			return ctx.JSON(http.StatusUnauthorized, noObjectEditorPrivilege)
		}
	*/
	if !s.Store.Exist(objectId) {
		msg := fmt.Sprintf("Update(%d) cannot modify this id, it does not exist !", objectId)
		s.Log.Error(msg)
		return ctx.JSON(http.StatusNotFound, msg)
	}
	tree := new(Tree)
	if err := ctx.Bind(tree); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Update has invalid format [%v]", err))
	}
	lastModificationUser := int32(claims.User.ExternalId)
	tree.LastModificationUser = &lastModificationUser
	if len(tree.Name) < 1 {
		return ctx.JSON(http.StatusBadRequest, "Tree name cannot be empty")
	}
	if len(tree.Name) < 5 {
		return ctx.JSON(http.StatusBadRequest, "Tree name minlength is 5")
	}
	if len(tree.Geom) < 1 {
		return ctx.JSON(http.StatusBadRequest, "Tree geom cannot be empty")
	}
	if (TreeAttributes{}) == tree.TreeAttributes {
		return ctx.JSON(http.StatusBadRequest, "Tree tree_attributes cannot be empty")
	}
	if tree.Id != objectId {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("Update id : [%d] and posted Id [%v] cannot differ ", objectId, tree.Id))
	}

	updatedTree, err := s.Store.Update(objectId, *tree)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Update got problem updating tree : %v", err))
	}
	return ctx.JSON(http.StatusOK, updatedTree)
}

func (s Service) GetMaxId(ctx echo.Context) error {
	s.Log.Debug("entering GetMaxId()")
	var maxTreeId int32 = 0
	maxTreeId, _ = s.Store.GetMaxId()
	s.Log.Info("# Exit GetMaxId() maxTreeId: %d", maxTreeId)
	return ctx.JSON(http.StatusOK, maxTreeId)
}

func (s Service) SearchTreesByName(ctx echo.Context, pattern string) error {
	s.Log.Debug("entering SearchTreesByName() pattern:%v\n", pattern)

	var search string = ""
	search = strings.TrimSpace(pattern)
	if search == "" || search == "*" {
		search = "%"
	}
	list, err := s.Store.SearchTreesByName(search)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("there was a problem when calling store.List :%v", err))
	}
	return ctx.JSON(http.StatusOK, list)

}

func (s *Service) ValidationList(ctx echo.Context, params ValidationListParams) error {
	s.Log.Debug("entering ValidationList() params:%v\n", params)

	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := claims.User.UserId
	s.Log.Info("in ValidationList : currentUserId: %d", currentUserId)
	if !(s.Store.IsObjectAdmin(int32(currentUserId)) || s.Store.IsObjectValidator(int32(currentUserId))) {
		return ctx.JSON(http.StatusUnauthorized, noObjectEditorPrivilege)
	}

	var sectName string = ""
	if params.Secteur != nil {
		sectName = string(*params.Secteur)
	}
	var idEmplacement int32 = -1
	if params.Emplacement != nil {
		idEmplacement = int32(*params.Emplacement)
	}
	list, err := s.Store.TreesToValidate(sectName, idEmplacement)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("there was a problem when calling store.TreesToValidate :%v", err))
	}
	return ctx.JSON(http.StatusOK, list)
}

func (s *Service) SaveValidation(ctx echo.Context) error {
	s.Log.Debug("entering SaveValidation()")

	claims := s.Server.JwtCheck.GetJwtCustomClaimsFromContext(ctx)
	currentUserId := claims.User.UserId
	externalUserId := claims.User.ExternalId
	s.Log.Info("in SaveValidation : external UserId: %d", externalUserId)
	if !(s.Store.IsObjectAdmin(int32(currentUserId)) || s.Store.IsObjectValidator(int32(currentUserId))) {
		return ctx.JSON(http.StatusUnauthorized, noObjectValidatorPrivilege)
	}

	var validationList []TreesToValidate
	if err := ctx.Bind(&validationList); err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("SaveValidation has invalid format [%v]", err))
	}
	if len(validationList) < 1 {
		return ctx.JSON(http.StatusBadRequest, "SaveValidation list cannot be empty")
	}
	for _, tree := range validationList {
		if tree.ExternalId < 1 {
			return ctx.JSON(http.StatusBadRequest, "SaveValidation tree externalId cannot be empty")
		}
		if tree.IsValidated {
			s.Store.ValidateTree(tree.ExternalId, tree.IsValidated, int32(externalUserId))
		}
	}

	return ctx.JSON(http.StatusOK, "SaveValidation")
}

func (s Service) GetDicoTable(ctx echo.Context, table GetDicoTableParamsTable) error {
	s.Log.Debug("entering GetDicoValidation(%s)", table)

	treedico, err := s.Store.GetDicoTable(table)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving dico(%s) :%v", table, err))
	}
	return ctx.JSON(http.StatusOK, treedico)

}

func (s Service) GetGestionComSecteurs(ctx echo.Context) error {
	s.Log.Debug("entering GetGestionComSecteurs")

	dico, err := s.Store.GetGestionComSecteurs()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving secteurs dico :%v", err))
	}
	return ctx.JSON(http.StatusOK, dico)
}

func (s Service) GetGestionComEmplacementsSecteur(ctx echo.Context, secteur string) error {
	s.Log.Debug("entering GetGestionComEmplacementsSecteur")

	dico, err := s.Store.GetGestionComEmplacementsSecteur(secteur)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving emplacements(%s) dico :%v", secteur, err))
	}
	return ctx.JSON(http.StatusOK, dico)
}

func (s Service) GetEmplacements(ctx echo.Context) error {
	s.Log.Debug("entering GetEmplacements")

	dico, err := s.Store.GetEmplacements()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving emplacements dico :%v", err))
	}
	return ctx.JSON(http.StatusOK, dico)
}

func (s Service) GetGestionComEmplacementsCentroidEmplacementId(ctx echo.Context, idEmplacement int32) error {
	s.Log.Debug("entering GetGestionComEmplacementsCentroidEmplacementId(%d)", idEmplacement)

	centroid, err := s.Store.GetGestionComEmplacementsCentroidEmplacementId(idEmplacement)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving centroid(%d) dico :%v", idEmplacement, err))
	}
	return ctx.JSON(http.StatusOK, centroid)
}

func (s Service) GetBuildingCenter(ctx echo.Context, addressId int32) error {
	s.Log.Debug("entering GetBuildingCenter(%d)", addressId)

	center, err := s.Store.GetBuildingCenter(addressId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving center(%d) :%v", addressId, err))
	}
	return ctx.JSON(http.StatusOK, center)
}

func (s Service) GetBuildingsNumbers(ctx echo.Context, streetId int32) error {
	s.Log.Debug("entering GetBuildingsNumbers(%d)", streetId)

	dico, err := s.Store.GetBuildingsNumbers(streetId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving buildingd numbers(%d) :%v", streetId, err))
	}
	return ctx.JSON(http.StatusOK, dico)
}

func (s Service) GetStreets(ctx echo.Context) error {
	s.Log.Debug("entering GetStreets")

	dico, err := s.Store.GetStreets()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("problem retrieving streets dico :%v", err))
	}
	return ctx.JSON(http.StatusOK, dico)
}

func (s *Service) GetGroupByName(ctx echo.Context, name string) error {
	s.Log.Debug("entering GetGroupByName(%s)", name)
	
    group, err := s.Store.GetGroupByName(name)
    if err != nil {
        return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("error getting group by name: %v", err))
    }
    return ctx.JSON(http.StatusOK, group)
}