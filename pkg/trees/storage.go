package trees

import (
	"errors"
	"fmt"

	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
)

// Storage is an interface to different implementation of persistence for Trees
type Storage interface {
	// List returns the list of existing objects with the given offset and limit.
	List(offset, limit int) ([]*TreeList, error)
	// Get returns the object with the specified objects ID.
	Get(id int32) (*Tree, error)
	// GetMaxId returns the maximum value of objects id existing in store.
	GetMaxId() (int32, error)
	// Exist returns true only if a objects with the specified id exists in store.
	Exist(id int32) bool
	// Count returns the total number of objects.
	Count() (int32, error)
	// Create saves a new objects in the storage.
	Create(object Tree) (*Tree, error)
	// Update updates the objects with given ID in the storage.
	Update(id int32, object Tree) (*Tree, error)
	// Delete removes the objects with given ID from the storage.
	Delete(id int32) error
	// SearchTreesByName list of existing objects where the name contains the given search pattern or err if not found
	SearchTreesByName(pattern string) ([]*TreeList, error)
	// TreesToValidate returns a list of trees to be validated given sector and emplacement
	TreesToValidate(secteur string, idEmplacement int32) ([]*ValidationList, error)

	ValidateTree(id int32, isValidated bool, idValidator int32) error
	// IsTreeActive returns true if the object with the specified id has the is_active attribute set to true
	IsTreeActive(id int32) bool

	IsUserAdmin(id int32) bool

	IsObjectAdmin(id int32) bool

	IsObjectEditor(id int32) bool

	IsObjectValidator(id int32) bool

	GetDicoTable(table GetDicoTableParamsTable) ([]*TreeDico, error)

	GetGestionComSecteurs() ([]*Dico, error)

	GetEmplacements() ([]*Dico, error)

	GetGestionComEmplacementsCentroidEmplacementId(int32) (*EmplacementCentroid, error)

	GetGestionComEmplacementsSecteur(string) ([]*Dico, error)

	GetBuildingCenter(int32) (*Center, error)

	GetBuildingsNumbers(int32) ([]*Dico, error)

	GetStreets() ([]*Dico, error)

	GetGroupByName(string) ([]*Group, error)
}

func GetStorageInstance(dbDriver string, db database.DB, l golog.MyLogger) (Storage, error) {
	var store Storage
	switch dbDriver {
	case "pgx":
		pgConn, err := db.GetPGConn()
		if err != nil {
			return nil, err
		}
		store = PGX{
			con: pgConn,
			log: l,
		}

	default:
		return nil, errors.New("unsupported DB driver type")
	}
	return store, nil
}

func GetErrorF(errMsg string, err error) error {
	return fmt.Errorf("%s [%v]", errMsg, err)
}
