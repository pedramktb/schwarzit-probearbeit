package testData

import (
	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
	"gorm.io/gorm"
)

var (
	TestUserOldVersionID = uuid.New()

	TestUserLatestVersionID = uuid.New()
)

var TestUser = types.User{
	ID:              uuid.New(),
	VersionID:       TestUserLatestVersionID,
	IsLatestVersion: true,
	FirstName:       "test",
	LastName:        "user",
	Email:           "test@test.com",
	Phone:           "+49123456789",
	IsAdmin:         false,
}

func MigrateTestData(db *gorm.DB) {
	migrateUsers(db)
}

func migrateUsers(db *gorm.DB) {
	if err := db.Table("users").Create(map[string]any{
		"id": TestUser.ID,
	}).Error; err != nil {
		panic(errors.Wrap(err, "failed to create Test data"))
	}

	if err := db.Table("user_versions").Create(map[string]any{
		"id":         TestUserOldVersionID,
		"user_id":    TestUser.ID,
		"first_name": "old",
		"last_name":  "user",
		"email":      TestUser.Email,
		"phone":      TestUser.Phone,
		"is_admin":   true,
	}).Error; err != nil {
		panic(errors.Wrap(err, "failed to create Test data"))
	}

	if err := db.Table("user_versions").Create(map[string]any{
		"id":         TestUserLatestVersionID,
		"user_id":    TestUser.ID,
		"first_name": TestUser.FirstName,
		"last_name":  TestUser.LastName,
		"email":      TestUser.Email,
		"phone":      TestUser.Phone,
	}).Error; err != nil {
		panic(errors.Wrap(err, "failed to create Test data"))
	}
}
