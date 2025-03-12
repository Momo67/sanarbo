package trees

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
)

const (
	noRecords = "records not found"
)

var ErrNoRecordFound = errors.New(noRecords)

type PGX struct {
	con *pgxpool.Pool
	log golog.MyLogger
}


func (P PGX) List(offset, limit int) ([]*TreeList, error) {
	P.log.Debug("entering List(%d, %d)", offset, limit)
	var res []*TreeList

	err := pgxscan.Select(context.Background(), P.con, &res, treesList, limit, offset)
	if err != nil {
		P.log.Error("List pgxscan.Select unexpectedly failed, error : %v", err)
		return nil, err
	}
	if res == nil {
		P.log.Info("List returned no results ")
		return nil, errors.New(noRecords)
	}

	return res, nil
}

func (P PGX) Get(id int32) (*Tree, error) {
	P.log.Debug("entering Get(%d)", id)
	res := &Tree{}

	err := pgxscan.Get(context.Background(), P.con, res, treesGet, id)
	if err != nil {
		P.log.Error("Get(%d) pgxscan.Get unexpectedly failed, error : %v", id, err)
		return nil, err
	}
	if res == (&Tree{}) {
		P.log.Info("Get(%d) returned no results ", id)
		return nil, errors.New(noRecords)
	}
	return res, nil
}

func (P PGX) GetMaxId() (int32, error) {
	P.log.Debug("entering GetMaxId()")
	var existingMaxId int32
	err := P.con.QueryRow(context.Background(), treesGetMaxId).Scan(&existingMaxId)
	if err != nil {
		P.log.Error("GetMaxId() could not be retrieved from DB. failed QueryRow.Scan err: %v", err)
		return 0, err
	}
	return existingMaxId, nil
}

func (P PGX) Exist(id int32) bool {
	P.log.Debug("entering Exist(%d)", id)
	var count int32 = 0

	err := P.con.QueryRow(context.Background(), treesExist, id).Scan(&count)
	if err != nil {
		P.log.Error("Exist(%d) could not be retrieved from DB. failed QueryRow.Scan err: %v", id, err)
		return false
	}
	if count > 0 {
		P.log.Info("Exist(%d) id does exist count:%v", id, count)
		return true
	} else {
		P.log.Info("Exist(%d) id does not exist count:%v", id, count)
		return false
	}
}

func (P PGX) Count() (int32, error) {
	P.log.Debug("entering Count()")
	var count int32
	err := P.con.QueryRow(context.Background(), treesCount).Scan(&count)
	if err != nil {
		P.log.Error("Count() could not be retrieved from DB. failed Query.Scan err: %v", err)
		return 0, err
	}
	return count, nil
}

func (P PGX) Create(object Tree) (*Tree, error) {
	P.log.Debug("entering Create(%q,%q,%#v)", object.Name, object.Geom, object.TreeAttributes)
	var lastInsertId int = 0

	err := P.con.QueryRow(context.Background(), treesCreate,
		object.Name, &object.Description, object.ExternalId, object.IsActive, &object.Comment, object.Creator, object.Geom, object.TreeAttributes).Scan(&lastInsertId)
	if err != nil {
		P.log.Error("Create(%q) unexpectedly failed. error : %v", object.Name, err)
		return nil, err
	}
	P.log.Info("Create(%q) created with id : %v", object.Name, lastInsertId)

	createdTree, err := P.Get(int32(lastInsertId))
	if err != nil {
		return nil, GetErrorF("error : tree was created, but cannot be retrieved", err)
	}
	return createdTree, nil
}

func (P PGX) Update(id int32, object Tree) (*Tree, error) {
	P.log.Debug("entering Update(%q,%q,%#v)", object.Name, object.Geom, object.TreeAttributes)

	now := time.Now()
	object.LastModificationTime = &now
	if !object.IsActive {
		object.InactivationTime = &now
	} else {
		object.InactivationTime = nil
	}
	P.log.Info("Just before Update(%+v)", object)

	res, err := P.con.Exec(context.Background(), treesUpdate,
		object.Name, &object.Description, &object.ExternalId, object.IsActive, &object.InactivationTime, &object.InactivationReason,
		&object.Comment, &object.IsValidated, &object.IdValidator, &object.LastModificationUser, object.Geom, &object.TreeAttributes, id)
	if err != nil {
		return nil, GetErrorF("error : Update() query failed", err)
	}
	if res.RowsAffected() < 1 {
		return nil, GetErrorF("error : Update() no row modified", err)
	}
	updatedTree, err := P.Get(id)
	if err != nil {
		return nil, GetErrorF("error : Update() user updated, but cannot be retrieved", err)
	}
	return updatedTree, nil
}

func (P PGX) Delete(id int32) error {
	P.log.Debug("entering Delete(%d)", id)

	res, err := P.con.Exec(context.Background(), treesDelete, id)
	if err != nil {
		return GetErrorF("error : tree could not be deleted", err)
	}
	if res.RowsAffected() < 1 {
		return GetErrorF("error : tree was not deleted", err)
	}

	return nil
}

func (P PGX) SearchTreesByName(pattern string) ([]*TreeList, error) {
	P.log.Debug("entering SearchTreesByName(%s)", pattern)
	var res []*TreeList
	var search string = ""

	search = strings.TrimSpace(pattern)
	search = strings.ReplaceAll(search, "*", "%")
	err := pgxscan.Select(context.Background(), P.con, &res, treesSearchByName, search)
	if err != nil {
		return nil, GetErrorF("error : SearchTreesByName query failed", err)
	}
	if res == nil {
		P.log.Info("SearchTreesByName returned no results ")
		return nil, ErrNoRecordFound
	}

	return res, nil
}

// TreesToValidate implements Storage.
func (P PGX) TreesToValidate(sectName string, idEmplacement int32) ([]*ValidationList, error) {
	P.log.Debug("entering TreesToValidate(%s, %d)", sectName, idEmplacement)
	var res []*ValidationList

	var secteur *string
	if sectName == "" {
		secteur = nil
	} else {
		secteur = &sectName
	}
	var emplacement *int32
	if idEmplacement == -1 {
		emplacement = nil
	} else {
		emplacement = &idEmplacement
	}
	err := pgxscan.Select(context.Background(), P.con, &res, treesToValidate, &secteur, &emplacement)
	if err != nil {
		return nil, GetErrorF("error : TreesToValidate query failed", err)
	}
	if res == nil {
		P.log.Info("TreesToValidate returned no results ")
		return []*ValidationList{}, nil
	}

	return res, nil
}

func (P PGX) ValidateTree(id int32, isValidated bool, idValidator int32) error {
	P.log.Debug("entering ValidateTree(%d, %d, %v)", id, idValidator, isValidated)
	_, err := P.con.Exec(context.Background(), validateTrees, id, isValidated, idValidator)
	if err != nil {
		return GetErrorF("error : ValidateTree query failed", err)
	}

	return nil
}

func (P PGX) IsTreeActive(id int32) bool {
	P.log.Debug("entering IsTreeActive(%d)", id)
	var isActive bool
	err := P.con.QueryRow(context.Background(), "SELECT is_active FROM tree_mobile WHERE id = $1", id).Scan(&isActive)
	if err != nil {
		P.log.Error("IsTreeActive(%d) could not be retrieved from DB. failed QueryRow.Scan err: %v", id, err)
		return false
	}
	return isActive
}

func (P PGX) IsUserAdmin(id int32) bool {
	P.log.Debug("entering IsUserAdmin(%d)", id)
	var isLocked bool
	var isUserAdmin bool
	err := P.con.QueryRow(context.Background(), "SELECT is_locked, is_admin FROM go_user WHERE id = $1", id).Scan(&isLocked, &isUserAdmin)
	if err != nil {
		P.log.Error("IsUserAdmin(%d) could not be retrieved from DB. failed QueryRow.Scan err: %v", id, err)
		return false
	}
	if isLocked {
		P.log.Error("IsUserAdmin(%d) user is locked", id)
		return false
	}
	return isUserAdmin
}

func (P PGX) IsObjectAdmin(id int32) bool {
	P.log.Debug("entering IsObjectAdmin(%d)", id)
	var isLocked bool
	var isObjectAdmin bool
	err := P.con.QueryRow(context.Background(), "SELECT is_locked, (2 = ANY(groups_id)) AS is_object_admin FROM go_user WHERE id = $1", id).Scan(&isLocked, &isObjectAdmin)
	if err != nil {
		P.log.Error("IsObjectAdmin(%d) could not be retrieved from DB. failed QueryRow.Scan err: %v", id, err)
		return false
	}
	if isLocked {
		P.log.Error("IsObjectAdmin(%d) user is locked", id)
		return false
	}
	return isObjectAdmin
}

func (P PGX) IsObjectEditor(id int32) bool {
	P.log.Debug("entering IsObjectEditor(%d)", id)
	var isLocked bool
	var isObjectEditor bool
	err := P.con.QueryRow(context.Background(), "SELECT is_locked, (3 = ANY(groups_id)) AS is_object_editor FROM go_user WHERE id = $1", id).Scan(&isLocked, &isObjectEditor)
	if err != nil {
		P.log.Error("IsObjectEditor(%d) could not be retrieved from DB. failed QueryRow.Scan err: %v", id, err)
		return false
	}
	if isLocked {
		P.log.Error("IsObjectEditor(%d) user is locked", id)
		return false
	}
	return isObjectEditor
}

// IsObjectValidator implements Storage.
func (P PGX) IsObjectValidator(id int32) bool {
	P.log.Debug("entering IsObjectValidator(%d)", id)
	var isLocked bool
	var isObjectValidator bool
	err := P.con.QueryRow(context.Background(), "SELECT is_locked, (6 = ANY(groups_id)) AS is_object_validator FROM go_user WHERE id = $1", id).Scan(&isLocked, &isObjectValidator)
	if err != nil {
		P.log.Error("IsObjectValidator(%d) could not be retrieved from DB. failed QueryRow.Scan err: %v", id, err)
		return false
	}
	if isLocked {
		P.log.Error("IsObjectValidator(%d) user is locked", id)
		return false
	}
	return isObjectValidator
}

func CompareTree(t1, t2 *TreeList, attr string) bool {
	if t1 == t2 {
		return true
	}

	inrec, _ := json.Marshal(t1)
	var m1 map[string]interface{}
	json.Unmarshal(inrec, &m1)
	var m2 map[string]interface{}
	inrec, _ = json.Marshal(t2)
	json.Unmarshal(inrec, &m2)
	return m1[attr] == m2[attr]
}

func (P PGX) GetDicoTable(table GetDicoTableParamsTable) ([]*TreeDico, error) {
	P.log.Debug("entering GetDico(%s)", table)
	var res []*TreeDico
	var query string = ""

	switch table {
	case Validation:
		query = treesDicoGetValidation
	case ToBeChecked:
		query = treesDicoGetToBeChecked
	case Note:
		query = treesDicoGetNote
	case Entourage:
		query = treesDicoGetEntourage
	case Check:
		query = treesDicoGetChk
	case RevSurface:
		query = treesDicoGetRevSurface
	case EtatSanitaire:
		query = treesDicoGetEtatSanitaire
	case EtatSanitaireRem:
		query = treesDicoGetEtatSanitaireRem
	default:
		return nil, errors.New("error : GetDico table unknown")
	}

	if query != "" {
		err := pgxscan.Select(context.Background(), P.con, &res, query)
		if err != nil {
			P.log.Error("GetDico pgxscan.Select unexpectedly failed, error : %v", err)
			return nil, err
		}
		if res == nil {
			P.log.Info("GetDico returned no results ")
			return nil, errors.New(noRecords)
		}

		return res, nil
	} else {
		return nil, errors.New("error : GetDico table not specified")
	}
}

func (P PGX) GetGestionComSecteurs() ([]*Dico, error) {
	P.log.Debug("entering GetGestionComSecteurs")
	var res []*Dico

	err := pgxscan.Select(context.Background(), P.con, &res, secteursList)
	if err != nil {
		P.log.Error("GetGestionComSecteurs pgxscan.Select unexpectedly failed, error : %v", err)
		return nil, err
	}
	if res == nil {
		P.log.Info("GetGestionComSecteurs returned no results ")
		return nil, errors.New(noRecords)
	}

	return res, nil
}

func (P PGX) GetEmplacements() ([]*Dico, error) {
	P.log.Debug("entering GetEmplacements")
	var res []*Dico

	err := pgxscan.Select(context.Background(), P.con, &res, emplacementsList)
	if err != nil {
		P.log.Error("GetEmplacements pgxscan.Select unexpectedly failed, error : %v", err)
		return nil, err
	}
	if res == nil {
		P.log.Info("GetEmplacements returned no results ")
		return nil, errors.New(noRecords)
	}

	return res, nil
}

func (P PGX) GetGestionComEmplacementsSecteur(secteur string) ([]*Dico, error) {
	P.log.Debug("entering GetGestionComEmplacementsSecteur(%s)", secteur)
	var res []*Dico

	err := pgxscan.Select(context.Background(), P.con, &res, emplacementsListBySecteur, secteur)
	if err != nil {
		P.log.Error("GetGestionComEmplacementsSecteur(%s) pgxscan.Select unexpectedly failed, error : %v", secteur, err)
		return nil, err
	}
	if res == nil {
		P.log.Info("GetGestionComEmplacementsSecteur returned no results ")
		return nil, errors.New(noRecords)
	}

	return res, nil
}

func (P PGX) GetGestionComEmplacementsCentroidEmplacementId(idemplacement int32) (*EmplacementCentroid, error) {
	P.log.Debug("entering GetGestionComEmplacementsCentroidEmplacementId(%d)", idemplacement)

	res := &EmplacementCentroid{}
	err := pgxscan.Get(context.Background(), P.con, res, emplacementCentroid, idemplacement)
	if err != nil {
		P.log.Error("GetGestionComEmplacementsCentroidEmplacementId(%d) pgxscan.Get unexpectedly failed, error : %v", idemplacement, err)
		return nil, err
	}
	if res == (&EmplacementCentroid{}) {
		P.log.Info("GetGestionComEmplacementsCentroidEmplacementId(%d) returned no results ", idemplacement)
		return nil, errors.New(noRecords)
	}
	return res, nil
}

func (P PGX) GetBuildingCenter(idaddress int32) (*Center, error) {
	P.log.Debug("entering GetBuildingCenter(%d)", idaddress)

	res := &Center{}
	err := pgxscan.Get(context.Background(), P.con, res, buildingCenter, idaddress)
	if err != nil {
		P.log.Error("GetBuildingCenter(%d) pgxscan.Get unexpectedly failed, error : %v", idaddress, err)
		return nil, err
	}
	if res == (&Center{}) {
		P.log.Info("GetBuildingCenter(%d) returned no results ", idaddress)
		return nil, errors.New(noRecords)
	}
	return res, nil
}

func (P PGX) GetBuildingsNumbers(idstreet int32) ([]*Dico, error) {
	P.log.Debug("entering GetBuildingsNumbers(%d)", idstreet)

	var res []*Dico
	err := pgxscan.Select(context.Background(), P.con, &res, buildingsNumberByStreet, idstreet)
	if err != nil {
		P.log.Error("GetBuildingsNumbers(%d) pgxscan.Select unexpectedly failed, error : %v", idstreet, err)
		return nil, err
	}
	if res == nil {
		P.log.Info("Select(%d) returned no results ", idstreet)
		return nil, errors.New(noRecords)
	}
	return res, nil
}

func (P PGX) GetStreets() ([]*Dico, error) {
	P.log.Debug("entering GetStreets")

	var res []*Dico
	err := pgxscan.Select(context.Background(), P.con, &res, streetsList)
	if err != nil {
		P.log.Error("GetStreets pgxscan.Select unexpectedly failed, error : %v", err)
		return nil, err
	}
	if res == nil {
		P.log.Info("Select returned no results ")
		return nil, errors.New(noRecords)
	}
	return res, nil
}

func (P PGX) GetGroupByName(name string) ([]*Group, error) {
	P.log.Debug("entering GetGroupByName(%s)", name)

	var res []*Group
	err := pgxscan.Select(context.Background(), P.con, &res, groupByName, name)
	if err != nil {
		P.log.Error("GetGroupByName(%s) pgxscan.Select unexpectedly failed, error : %v", name, err)
		return nil, err	
	}	
	if res == nil {
		P.log.Info("GetGroupByName(%s) returned no results ", name)
		return nil, errors.New(noRecords)
	}
	return res, nil
}
