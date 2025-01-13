package userDB

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
)

type db struct {
	*gorm.DB
}

func create(g *gorm.DB) *db {
	return &db{
		DB: g,
	}
}

func lastVersionQuery(db, tx *gorm.DB) *gorm.DB {
	return tx.Joins("JOIN (?) AS last_version ON users.id = last_version.user_id",
		db.Table("user_versions").Select("DISTINCT ON (user_id) *").Order("user_id, created_at DESC"),
	).Select(
		"users.id as id",
		"last_version.id as version_id",
		"last_version.first_name as first_name",
		"last_version.last_name as last_name",
		"last_version.email as email",
		"last_version.phone as phone",
		"last_version.is_admin as is_admin",
		"last_version.password_hash as password_hash",
		"true as is_latest_version",
	).Where("users.deleted_at IS NULL")
}

func allVersionsQuery(db, tx *gorm.DB) *gorm.DB {
	return tx.Table("users").Joins("JOIN user_versions ON users.id = user_versions.user_id").Joins(
		"JOIN (?) AS last_version ON users.id = last_version.user_id",
		db.Table("user_versions").Select("DISTINCT ON (user_id) user_id", "id").Order("user_id, created_at DESC"),
	).Select(
		"users.id as id",
		"user_versions.id as version_id",
		"user_versions.first_name as first_name",
		"user_versions.last_name as last_name",
		"user_versions.email as email",
		"user_versions.phone as phone",
		"user_versions.is_admin as is_admin",
		"user_versions.password_hash as password_hash",
		"user_versions.id = last_version.id as is_latest_version",
	).Where("users.deleted_at IS NULL")
}

func (d *db) Get(ctx context.Context, id uuid.UUID) (types.User, error) {
	var user types.User
	err := lastVersionQuery(d.WithContext(ctx), d.WithContext(ctx).Table("users")).
		Where("users.id = ?", id).First(&user).Error
	return user, types.DBError(err)
}

func (d *db) Query(ctx context.Context, params types.QueryParams) ([]types.User, error) {
	var users []types.User
	err := types.Query(lastVersionQuery(d.WithContext(ctx), d.WithContext(ctx).Table("users")), params).Find(&users).Error
	return users, types.DBError(err)
}

func (d *db) Save(ctx context.Context, user types.User) (types.User, error) {
	base, version := user.ToSave()
	err := d.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create or find the base user
		if base != nil {
			if err := tx.Table("users").Create(base).Error; err != nil {
				return types.DBError(err)
			}
		} else if err := tx.Table("users").Where("id = ? AND deleted_at IS NULL", user.ID).First(&types.User{}).Error; err != nil {
			return types.DBError(err)
		}

		// Create the new version
		if err := tx.Table("user_versions").Create(version).Error; err != nil {
			return types.DBError(err)
		}

		return nil
	})
	return user, err
}

func (d *db) Delete(ctx context.Context, id uuid.UUID) error {
	if err := d.WithContext(ctx).Table("users").Where("id = ? AND deleted_at IS NULL", id).First(&types.User{}).Error; err != nil {
		return types.DBError(err)
	}
	return types.DBError(d.WithContext(ctx).Table("users").Where("id = ? AND deleted_at IS NULL", id).Delete(nil).Error)
}

func (d *db) GetVersion(ctx context.Context, versionID uuid.UUID) (types.User, error) {
	var user types.User
	err := allVersionsQuery(d.WithContext(ctx), d.WithContext(ctx).Table("users")).
		Where("user_versions.id = ?", versionID).First(&user).Error
	return user, types.DBError(err)
}

func (d *db) GetByEmail(ctx context.Context, email string) (types.User, error) {
	var user types.User
	err := lastVersionQuery(d.WithContext(ctx), d.WithContext(ctx).Table("users")).
		Where("email = ?", email).First(&user).Error
	return user, types.DBError(err)
}
