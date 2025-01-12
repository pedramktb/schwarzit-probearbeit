package userDB

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/google/uuid"
	testData "github.com/pedramktb/schwarzit-probearbeit/internal/test_data"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
	"github.com/pedramktb/schwarzit-probearbeit/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

var postgresContainer testcontainers.Container
var ip, port string

func TestMain(m *testing.M) {
	postgresContainer, ip, port = postgres.Test_Create_Container()
	defer func(postgresContainer testcontainers.Container, ctx context.Context) {
		_ = postgresContainer.Terminate(ctx)
	}(postgresContainer, context.Background())

	defer os.Exit(m.Run())
}

func Test_Get(t *testing.T) {
	dbName := "test-user-get"
	db := postgres.Test_Create_DB(ip, port, dbName)
	defer postgres.Test_Drop_DB(db, ip, port, dbName)
	testData.MigrateTestData(db)

	// test
	tests := []struct {
		name    string
		id      uuid.UUID
		want    types.User
		wantErr bool
	}{
		{
			name:    "Success Case",
			id:      testData.TestUser.ID,
			want:    testData.TestUser,
			wantErr: false,
		},
		{
			name:    "Not Found Case",
			id:      uuid.New(),
			want:    types.User{},
			wantErr: true,
		},
	}

	userDB := create(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userDB.Get(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Query(t *testing.T) {
	dbName := "test-user-query"
	db := postgres.Test_Create_DB(ip, port, dbName)
	defer postgres.Test_Drop_DB(db, ip, port, dbName)
	testData.MigrateTestData(db)

	// test
	tests := []struct {
		name    string
		query   types.QueryParams
		want    []types.User
		wantErr bool
	}{
		{
			name:    "Success Case",
			query:   types.QueryParams{Conditions: &types.UserPatch{FirstName: types.ToOptional("test")}},
			want:    []types.User{testData.TestUser},
			wantErr: false,
		},
		{
			name:    "Not Found Case",
			query:   types.QueryParams{Conditions: &types.UserPatch{FirstName: types.ToOptional("not found")}},
			want:    []types.User{},
			wantErr: false,
		},
	}

	userDB := create(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userDB.Query(context.Background(), tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			assert.Equal(t, len(tt.want), len(got))
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func Test_Save(t *testing.T) {
	dbName := "test-user-save"
	db := postgres.Test_Create_DB(ip, port, dbName)
	defer postgres.Test_Drop_DB(db, ip, port, dbName)
	testData.MigrateTestData(db)

	TestUser1 := testData.TestUser
	TestUser1.LastName = "user 1"

	TestUser2 := testData.TestUser
	TestUser2.ID = uuid.Nil
	TestUser2.LastName = "user 2"

	// test
	tests := []struct {
		name    string
		user    types.User
		want    types.User
		wantErr bool
	}{
		{
			name:    "Update Case",
			user:    TestUser1,
			want:    TestUser1,
			wantErr: false,
		},
		{
			name:    "New Case",
			user:    TestUser2,
			want:    TestUser2,
			wantErr: false,
		},
	}

	userDB := create(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saved, err := userDB.Save(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			if tt.want.ID == uuid.Nil {
				tt.want.ID = saved.ID
			}
			tt.want.VersionID = saved.VersionID
			assert.Equal(t, tt.want.ID, saved.ID)
			got, err := userDB.Get(context.Background(), saved.ID)
			if err != nil {
				t.Errorf("db.Get() error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Delete(t *testing.T) {
	dbName := "test-user-delete"
	db := postgres.Test_Create_DB(ip, port, dbName)
	defer postgres.Test_Drop_DB(db, ip, port, dbName)
	testData.MigrateTestData(db)

	// test
	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "Success Case",
			id:      testData.TestUser.ID,
			wantErr: false,
		},
		{
			name:    "Not Found Case",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	userDB := create(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userDB.Delete(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			got, err := userDB.Get(context.Background(), tt.id)
			if errors.Is(err, types.ErrNotFound) {
			} else if err != nil {
				t.Errorf("db.Get() error = %v", err)
				return
			}
			assert.Equal(t, types.User{}, got)
		})
	}
}

func Test_GetVersion(t *testing.T) {
	dbName := "test-user-get-version"
	db := postgres.Test_Create_DB(ip, port, dbName)
	defer postgres.Test_Drop_DB(db, ip, port, dbName)
	testData.MigrateTestData(db)

	wantOldVersion := testData.TestUser
	wantOldVersion.VersionID = testData.TestUserOldVersionID
	wantOldVersion.IsLatestVersion = false
	wantOldVersion.FirstName = "old"

	// test
	tests := []struct {
		name      string
		versionID uuid.UUID
		want      types.User
		wantErr   bool
	}{
		{
			name:      "Success Case",
			versionID: testData.TestUser.VersionID,
			want:      testData.TestUser,
			wantErr:   false,
		},
		{
			name:      "Success Case Old Version",
			versionID: testData.TestUserOldVersionID,
			want:      wantOldVersion,
			wantErr:   false,
		},
		{
			name:      "Not Found Case",
			versionID: uuid.New(),
			want:      types.User{},
			wantErr:   true,
		},
	}

	userDB := create(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userDB.GetVersion(context.Background(), tt.versionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.GetVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetByEmail(t *testing.T) {
	dbName := "test-user-get-by-email"
	db := postgres.Test_Create_DB(ip, port, dbName)
	defer postgres.Test_Drop_DB(db, ip, port, dbName)
	testData.MigrateTestData(db)

	// test
	tests := []struct {
		name    string
		email   string
		want    types.User
		wantErr bool
	}{
		{
			name:    "Success Case",
			email:   testData.TestUser.Email,
			want:    testData.TestUser,
			wantErr: false,
		},
		{
			name:    "Not Found Case",
			email:   "not@fou.nd",
			want:    types.User{},
			wantErr: true,
		},
	}

	userDB := create(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userDB.GetByEmail(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("db.GetByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
