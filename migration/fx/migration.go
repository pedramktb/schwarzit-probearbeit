package migrationDI

import (
	"context"

	"github.com/pedramktb/schwarzit-probearbeit/migration"
)

type migrator struct {
	migrators []migration.Migrator
}

func create(migrators ...migration.Migrator) *migrator {
	return &migrator{
		migrators: migrators,
	}
}

func (m *migrator) Migrate(ctx context.Context) {
	for _, migrator := range m.migrators {
		migrator.Migrate(ctx)
	}
}
