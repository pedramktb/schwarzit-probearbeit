package testData

import (
	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
	"gorm.io/gorm"
)

var (
	TestUserOldVersionID = uuid.New()
)

var TestUser = types.User{
	ID:              uuid.New(),
	VersionID:       uuid.New(),
	IsLatestVersion: true,
	FirstName:       "test",
	LastName:        "user",
	Email:           "test@test.com",
	Phone:           "+49123456789",
	IsAdmin:         false,
	PasswordHash:    "password",
}

var TestAdminUser = types.User{
	ID:              uuid.New(),
	VersionID:       uuid.New(),
	IsLatestVersion: true,
	FirstName:       "test",
	LastName:        "admin",
	Email:           "admin@test.com",
	Phone:           "+49123456789",
	IsAdmin:         true,
	PasswordHash:    "password",
}

func MigrateTestData(db *gorm.DB) {
	migrateUsers(db)
}

func migrateUsers(db *gorm.DB) {
	if err := db.Table("users").Create([]map[string]any{
		{
			"id": TestUser.ID,
		},
		{
			"id": TestAdminUser.ID,
		},
	}).Error; err != nil {
		panic(errors.Wrap(err, "failed to create Test data"))
	}

	if err := db.Table("user_versions").Create([]map[string]any{
		{
			"id":            TestUserOldVersionID,
			"user_id":       TestUser.ID,
			"first_name":    "old",
			"last_name":     "user",
			"email":         TestUser.Email,
			"phone":         TestUser.Phone,
			"is_admin":      TestUser.IsAdmin,
			"password_hash": TestUser.PasswordHash,
		},
	}).Error; err != nil {
		panic(errors.Wrap(err, "failed to create Test data"))
	}

	if err := db.Table("user_versions").Create([]map[string]any{
		{
			"id":            TestUser.VersionID,
			"user_id":       TestUser.ID,
			"first_name":    TestUser.FirstName,
			"last_name":     TestUser.LastName,
			"email":         TestUser.Email,
			"phone":         TestUser.Phone,
			"is_admin":      TestUser.IsAdmin,
			"password_hash": TestUser.PasswordHash,
		},
		{
			"id":            TestAdminUser.VersionID,
			"user_id":       TestAdminUser.ID,
			"first_name":    TestAdminUser.FirstName,
			"last_name":     TestAdminUser.LastName,
			"email":         TestAdminUser.Email,
			"phone":         TestAdminUser.Phone,
			"is_admin":      TestAdminUser.IsAdmin,
			"password_hash": TestAdminUser.PasswordHash,
		},
	}).Error; err != nil {
		panic(errors.Wrap(err, "failed to create Test data"))
	}
}
