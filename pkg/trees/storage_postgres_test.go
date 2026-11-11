package trees

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/sanarbo/pkg/version"
)

const (
	defaultDBPort    = 5432
	defaultDBIp      = "127.0.0.1"
	defaultDBSslMode = "prefer"
)

type WorkingEnv struct {
	l        golog.MyLogger
	logLevel golog.Level
	dbConn   database.DB
	storage  Storage
}

func TestSearchTreesByName(t *testing.T) {
	// this tree is inserted by the migration 000002 and also exists in the Goeland data
	const culotteSuisse = "Pyrus communis 'Culotte suisse' - secteur OUEST - Valency plantage (45904)"

	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()

	tests := []struct {
		name         string
		text         string
		limit        int
		wantContains string
		wantEmpty    bool
	}{
		{name: "search is case insensitive", text: "culotte SUISSE' - secteur ouest", limit: 100, wantContains: culotteSuisse},
		{name: "search finds a text in the middle of the name", text: "ulotte suis", limit: 100, wantContains: culotteSuisse},
		{name: "% is not a wildcard", text: "Culotte%suisse", limit: 100, wantEmpty: true},
		{name: "_ is not a wildcard", text: "_ulotte suisse", limit: 100, wantEmpty: true},
		{name: "no tree found returns an empty list", text: "666 no tree has this name", limit: 100, wantEmpty: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotErr := w.storage.SearchTreesByName(tt.text, tt.limit)

			assertNoError(t, gotErr)
			if tt.wantEmpty {
				assertEmptyJsonList(t, gotRes)
				return
			}
			for _, tree := range gotRes {
				if tree.Name == tt.wantContains {
					return
				}
			}
			t.Errorf("SearchTreesByName(%q) returned %d trees, none named %q", tt.text, len(gotRes), tt.wantContains)
		})
	}
	t.Run("limit is applied", func(t *testing.T) {
		gotRes, gotErr := w.storage.SearchTreesByName("Pyrus", 1)
		assertNoError(t, gotErr)
		if len(gotRes) != 1 {
			t.Errorf("SearchTreesByName(\"Pyrus\", 1) returned %d trees, want 1", len(gotRes))
		}
	})
}

// TestValidateTreesIsAllOrNothing checks that ValidateTrees runs in a transaction :
// when one external id does not exist, the other trees of the list must not be validated.
// It does not modify any data unless ValidateTrees is broken.
func TestValidateTreesIsAllOrNothing(t *testing.T) {
	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()
	skipIfTablesMissing(t, w, "thing", "lien_thing_thing")

	toValidate, err := w.storage.TreesToValidate(-1, -1)
	assertNoError(t, err)
	if len(toValidate) == 0 {
		t.Skip("no tree waiting for validation in db, cannot test the rollback")
	}
	existingId := toValidate[0].ExternalId
	const unknownId = -1

	err = w.storage.ValidateTrees([]int32{existingId, unknownId}, 0)

	if !errors.Is(err, ErrNoRecordFound) {
		t.Fatalf("ValidateTrees with an unknown id should return ErrNoRecordFound, got: %v", err)
	}
	stillToValidate, err := w.storage.TreesToValidate(-1, -1)
	assertNoError(t, err)
	for _, tree := range stillToValidate {
		if tree.ExternalId == existingId {
			return
		}
	}
	t.Errorf("tree with external_id %d was validated although the transaction failed", existingId)
}

// TestGetUnknownTreeReturnsErrNoRecordFound checks that an id missing in db is reported as ErrNoRecordFound
// (so the handler can answer 404) and not as a generic db error.
func TestGetUnknownTreeReturnsErrNoRecordFound(t *testing.T) {
	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()

	tree, err := w.storage.Get(-1)

	if !errors.Is(err, ErrNoRecordFound) {
		t.Fatalf("Get(-1) should return ErrNoRecordFound, got tree: %v, err: %v", tree, err)
	}
}

func TestGetGroupByName(t *testing.T) {
	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()

	t.Run("unknown group returns ErrNoRecordFound", func(t *testing.T) {
		group, err := w.storage.GetGroupByName("sanarbo-test-no-such-group")
		if !errors.Is(err, ErrNoRecordFound) {
			t.Fatalf("GetGroupByName should return ErrNoRecordFound, got group: %v, err: %v", group, err)
		}
	})
	t.Run("group created by the migrations is found", func(t *testing.T) {
		group, err := w.storage.GetGroupByName("object_validator")
		assertNoError(t, err)
		if group.Name != "object_validator" || group.Id < 1 {
			t.Errorf("GetGroupByName(object_validator) returned %+v", group)
		}
	})
}

func TestGetBuildingCenterAndCentroidNotFound(t *testing.T) {
	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()
	skipIfTablesMissing(t, w, "thi_street_building_address", "geodata_gestion_com.spadom_surfaces")
	pool, err := w.dbConn.GetPGConn()
	assertNoError(t, err)

	t.Run("unknown address returns ErrNoRecordFound", func(t *testing.T) {
		center, err := w.storage.GetBuildingCenter(-1)
		if !errors.Is(err, ErrNoRecordFound) {
			t.Fatalf("GetBuildingCenter(-1) should return ErrNoRecordFound, got center: %v, err: %v", center, err)
		}
	})
	t.Run("existing address returns its center", func(t *testing.T) {
		var idAddress int32
		assertNoError(t, pool.QueryRow(context.Background(), "SELECT idaddress FROM thi_street_building_address LIMIT 1").Scan(&idAddress))
		center, err := w.storage.GetBuildingCenter(idAddress)
		assertNoError(t, err)
		if center.Geometry == "" {
			t.Errorf("GetBuildingCenter(%d) returned an empty geometry", idAddress)
		}
	})
	t.Run("emplacement without surface returns ErrNoRecordFound", func(t *testing.T) {
		centroid, err := w.storage.GetGestionComEmplacementsCentroidEmplacementId(-1)
		if !errors.Is(err, ErrNoRecordFound) {
			t.Fatalf("GetGestionComEmplacementsCentroidEmplacementId(-1) should return ErrNoRecordFound, got centroid: %v, err: %v", centroid, err)
		}
	})
	t.Run("existing emplacement returns its centroid", func(t *testing.T) {
		var idEmplacement int32
		assertNoError(t, pool.QueryRow(context.Background(), "SELECT idgo_empl FROM geodata_gestion_com.spadom_surfaces WHERE idgo_empl IS NOT NULL LIMIT 1").Scan(&idEmplacement))
		centroid, err := w.storage.GetGestionComEmplacementsCentroidEmplacementId(idEmplacement)
		assertNoError(t, err)
		if centroid.Geometry == "" || centroid.Surface <= 0 {
			t.Errorf("GetGestionComEmplacementsCentroidEmplacementId(%d) returned %+v", idEmplacement, centroid)
		}
	})
}

// TestListsWithoutResultAreEmpty checks that a query without result returns an empty list and no error,
// so the handlers answer 200 with [] (and not 500, nor null which the frontend cannot use).
func TestListsWithoutResultAreEmpty(t *testing.T) {
	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()

	t.Run("List with an offset past the last tree", func(t *testing.T) {
		res, err := w.storage.List(1000000000, 10)
		assertNoError(t, err)
		assertEmptyJsonList(t, res)
	})
	t.Run("GetBuildingsNumbers of an unknown street", func(t *testing.T) {
		skipIfTablesMissing(t, w, "thi_street_building_address")
		res, err := w.storage.GetBuildingsNumbers(-1)
		assertNoError(t, err)
		assertEmptyJsonList(t, res)
	})
	t.Run("GetGestionComEmplacementsSecteurId of an unknown secteur", func(t *testing.T) {
		skipIfTablesMissing(t, w, "thing", "lien_thing_thing")
		res, err := w.storage.GetGestionComEmplacementsSecteurId(-1)
		assertNoError(t, err)
		assertEmptyJsonList(t, res)
	})
}

func TestTreesToValidateBySecteur(t *testing.T) {
	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()
	skipIfTablesMissing(t, w, "thing", "lien_thing_thing")

	secteurs, err := w.storage.GetGestionComSecteurs()
	assertNoError(t, err)
	if len(secteurs) == 0 {
		t.Skip("no secteur in db")
	}
	all, err := w.storage.TreesToValidate(-1, -1)
	assertNoError(t, err)

	bySecteur, err := w.storage.TreesToValidate(secteurs[0].Id, -1)

	assertNoError(t, err)
	if len(bySecteur) > len(all) {
		t.Errorf("TreesToValidate(secteur %d) returned %d trees, more than without filter (%d)", secteurs[0].Id, len(bySecteur), len(all))
	}
}

// TestUpdateKeepsInactivationTime creates its own tree, so it is skipped when the db user cannot insert in tree_mobile
// (the local sanarbo role has no right on seq_tree_mobile_id), it runs in the CI database.
func TestUpdateKeepsInactivationTime(t *testing.T) {
	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()

	circonference := int32(42)
	created, err := w.storage.Create(Tree{
		Name:           "Test inactivation time",
		IsActive:       true,
		Creator:        1,
		Geom:           "POINT(2538000 1152000)",
		TreeAttributes: TreeAttributes{Circonference: &circonference},
	})
	if err != nil && strings.Contains(err.Error(), "permission denied") {
		t.Skipf("the db user cannot create a tree in this database: %v", err)
	}
	assertNoError(t, err)
	defer func() { assertNoError(t, w.storage.Delete(created.Id)) }()

	update := func(isActive bool) *Tree {
		t.Helper()
		tree := *created
		tree.IsActive = isActive
		// sent by the client, must be ignored
		clientTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		tree.InactivationTime = &clientTime
		updated, err := w.storage.Update(created.Id, tree)
		assertNoError(t, err)
		return updated
	}

	firstInactivation := update(false).InactivationTime
	if firstInactivation == nil || firstInactivation.Year() == 2000 {
		t.Fatalf("inactivation_time should be set by the server when the tree becomes inactive, got %v", firstInactivation)
	}
	if again := update(false).InactivationTime; again == nil || !again.Equal(*firstInactivation) {
		t.Errorf("updating an inactive tree should keep inactivation_time %v, got %v", firstInactivation, again)
	}
	if reactivated := update(true).InactivationTime; reactivated != nil {
		t.Errorf("an active tree should have no inactivation_time, got %v", reactivated)
	}
	if inactivatedAgain := update(false).InactivationTime; inactivatedAgain == nil || inactivatedAgain.Before(*firstInactivation) {
		t.Errorf("inactivating the tree again should set a new inactivation_time after %v, got %v", firstInactivation, inactivatedAgain)
	}
}

func TestUserPrivilegesByGroupName(t *testing.T) {
	w := WorkingEnv{}
	assertWorkingEnvOK(t, w.Init())
	defer w.dbConn.Close()
	pool, err := w.dbConn.GetPGConn()
	assertNoError(t, err)

	findUser := func(t *testing.T, groupName string, member bool) int32 {
		t.Helper()
		var id int32
		err := pool.QueryRow(context.Background(), `
			SELECT u.id FROM go_user u
			WHERE NOT u.is_locked AND $2 = EXISTS (
				SELECT 1 FROM go_group g WHERE g.name = $1 AND g.is_active AND g.id = ANY(u.groups_id))
			ORDER BY u.id LIMIT 1`, groupName, member).Scan(&id)
		if err != nil {
			t.Skipf("no unlocked user with member=%v of group %s in this database", member, groupName)
		}
		return id
	}
	privileges := map[string]func(int32) bool{
		objectAdminGroup:     w.storage.IsObjectAdmin,
		objectEditorGroup:    w.storage.IsObjectEditor,
		objectValidatorGroup: w.storage.IsObjectValidator,
	}
	for groupName, hasPrivilege := range privileges {
		t.Run(groupName+" member", func(t *testing.T) {
			if id := findUser(t, groupName, true); !hasPrivilege(id) {
				t.Errorf("user %d is in group %s but has not the privilege", id, groupName)
			}
		})
		t.Run(groupName+" not member", func(t *testing.T) {
			if id := findUser(t, groupName, false); hasPrivilege(id) {
				t.Errorf("user %d is not in group %s but has the privilege", id, groupName)
			}
		})
		t.Run(groupName+" unknown user", func(t *testing.T) {
			if hasPrivilege(-1) {
				t.Errorf("unknown user has the privilege of group %s", groupName)
			}
		})
	}
}

func (w *WorkingEnv) Init() error {
	if w.logLevel == 0 {
		w.logLevel = golog.DebugLevel
	}

	err := w.GetLogger()
	if err != nil {
		return fmt.Errorf("got error when getting logger, err: %v", err)
	}
	err = w.GetDb()
	if err != nil {
		return fmt.Errorf("got error when getting db, err: %v", err)
	}
	err = w.GetStorage()
	if err != nil {
		return fmt.Errorf("got error when getting storage, err: %v", err)
	}

	return nil
}

func (w *WorkingEnv) GetLogger() error {
	log, err := golog.NewLogger("zap", w.logLevel, fmt.Sprintf("%s ", version.APP))
	if err != nil {
		return fmt.Errorf("got no logger error: %v", err)
	} else {
		w.l = log
	}
	return nil
}

func (w *WorkingEnv) GetDb() error {
	dbDsn := config.GetPgDbDsnUrlFromEnvOrPanic(defaultDBIp, defaultDBPort, tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	var dbConn database.DB
	dbConn, err := database.GetInstance("pgx", dbDsn, runtime.NumCPU(), w.l)
	if err != nil {
		return fmt.Errorf("error doing database.GetInstance(\"pgx\", dbDsn)  : %v\n", err)
	} else {
		w.dbConn = dbConn
	}
	return nil
}

func (w *WorkingEnv) GetStorage() error {
	treesStorage, err := GetStorageInstance("pgx", w.dbConn, w.l)
	if err != nil {
		return fmt.Errorf("error doing GetStorageInstance(\"pgx\", %#v, %#v) error: %v\n", w.dbConn, w.l, err)
	} else {
		w.storage = treesStorage
	}
	return nil
}

// skipIfTablesMissing skips the test when one of the tables does not exist, like the Goeland tables
// (thing, lien_thing_thing...) which are not created by the migrations, e.g. in the CI database.
func skipIfTablesMissing(t testing.TB, w WorkingEnv, tables ...string) {
	t.Helper()
	pool, err := w.dbConn.GetPGConn()
	if err != nil {
		t.Fatalf("cannot get db pool: %v", err)
	}
	for _, table := range tables {
		var exists bool
		err := pool.QueryRow(context.Background(), "SELECT to_regclass($1) IS NOT NULL", table).Scan(&exists)
		if err != nil {
			t.Fatalf("cannot check if table %s exists: %v", table, err)
		}
		if !exists {
			t.Skipf("table %s does not exist in this database (not created by the migrations)", table)
		}
	}
}

// assertEmptyJsonList checks that list is empty and is encoded as [] in json, not null
func assertEmptyJsonList(t testing.TB, list any) {
	t.Helper()
	got, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("cannot marshal list: %v", err)
	}
	if string(got) != "[]" {
		t.Errorf("got %s, want an empty json list []", got)
	}
}

func assertWorkingEnvOK(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Error(got.Error())
		t.Fatal("unable to get a working env")
	}
}

func assertName(t testing.TB, list TreeList, want string) {
	t.Helper()
	got := list.Name

	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}

func assertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one, got: %s", got)
	}
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
