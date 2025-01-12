package v1Migration

import (
	"context"
	_ "embed"

	"gorm.io/gorm"
)

type migrator struct {
	dst *gorm.DB
}

func create(dst *gorm.DB) *migrator {
	return &migrator{
		dst: dst,
	}
}

//go:embed migration.sql
var sqlMigration string

func (m *migrator) Migrate(ctx context.Context) {
	err := m.dst.WithContext(ctx).Exec(sqlMigration).Error
	if err != nil {
		panic(err)
	}
}
