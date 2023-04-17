package trees

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strings"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/sanarbo/pkg/version"
)

const (
	defaultDBPort              = 5432
	defaultDBIp                = "127.0.0.1"
	defaultDBSslMode           = "prefer"
)

func TestSearchTreesByName(t *testing.T) {
	type args struct {
		t string
	}

	prefix := fmt.Sprintf("%s ", version.APP)
	l, err := golog.NewLogger("zap", golog.DebugLevel, prefix)
	if err != nil {
		log.Fatalf("💥💥 error log.NewLogger error: %v\n", err)
	}
	/*
	secret, err := config.GetJwtSecretFromEnv()
	if err != nil {
		l.Fatal("💥💥 error doing config.GetJwtSecretFromEnv() error: %v\n", err)
	}
	tokenDuration, err := config.GetJwtDurationFromEnv(60)
	if err != nil {
		l.Fatal("💥💥 error doing config.GetJwtDurationFromEnv() error: %v\n", err)
	}
	*/
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
		tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	if err != nil {
		l.Fatal("💥💥 error doing config.GetPgDbDsnUrlFromEnv error: %v\n", err)
	}
	var dbConn database.DB
	dbConn, err = database.GetInstance("pgx", dbDsn, runtime.NumCPU(), l)
	if err != nil {
		l.Fatal("💥💥 error doing database.GetInstance(\"pgx\", dbDsn)  : %v\n", err)
	}
	defer dbConn.Close()

	treesStorage, err := GetStorageInstance("pgx", dbConn, l)
	if err != nil {
		l.Fatal("💥💥 error doing GetStorageInstance(\"pgx\", %#v, %#v) error: %v\n", dbConn, l, err)
	}

	tests := []struct {
		name		string
		args		args
		testFunc	func(any, any) []*TreeList
		wantRes		[]*TreeList
		wantErr		error
	}{
		{
			name: "it should return an object with name attribute matching pattern containing *",
			args: 	args{t: "*Tre*"},
			testFunc: func(val any, test any) []*TreeList {
				var matched bool = false
				var toFind string = ""

				toFind = strings.ReplaceAll(test.(string), "*", "")
				toFind = strings.ReplaceAll(toFind, "%", "")
				for _, v := range val.([]*TreeList) {
					matched, _ = regexp.MatchString(toFind, v.Name)
				}
				if matched {
					return val.([]*TreeList)
				}
				return nil
			},
			wantRes: []*TreeList{{Name: "MyNewTree"}},
			wantErr: nil,
		},
		{
			name: "it should return an object with name attribute matching pattern containing %",
			args: 	args{t: "%%Tre%%"},
			testFunc: func(val any, test any) []*TreeList {
				var matched bool = false
				var toFind string = ""

				toFind = strings.ReplaceAll(test.(string), "*", "")
				toFind = strings.ReplaceAll(toFind, "%", "")
				for _, v := range val.([]*TreeList) {
					matched, _ = regexp.MatchString(toFind, v.Name)
				}
				if matched {
					return val.([]*TreeList)
				}
				return nil
			},
			wantRes: []*TreeList{{Name: "MyNewTree"}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, gotErr := treesStorage.SearchTreesByName(tt.args.t)
			// fmt.Printf("##GetType(%v)[%T] returns: (%v, error: %v), wants: (%v, error:%v)", tt.args.t, tt.args.t, gotRes, gotErr, tt.wantRes, tt.wantErr)
			if (tt.testFunc != nil) {
				if !CompareTree((tt.testFunc(gotRes, tt.args.t))[0], (tt.wantRes)[0], "Name") || gotErr != tt.wantErr {
					t.Errorf("SearchTreesByName(%v)[%T] got: (%v, error: %v), wants: (%v, error:%v)",
						tt.args.t, tt.args.t, gotRes, gotErr, tt.wantRes, tt.wantErr)
				}
			}
		})
	}
}
