package trees

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
)

const (
	noRecords = "records not found"
)

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
		P.log.Error("Get(%d) pgxscan.Select unexpectedly failed, error : %v", id, err)
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
		return nil, errors.New(noRecords)
	}

	return res, nil
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
	//TODO implement a better user check...
	//Now only user with id(999) (bill board) is considered as admin
	if id == 999 {
		return true
	} else {
		return false
	}
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
