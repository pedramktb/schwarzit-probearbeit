package main

import (
	"context"

	_ "github.com/pedramktb/go-base-lib/pkg/env"
	"github.com/pedramktb/schwarzit-probearbeit/migration"
	migrationDI "github.com/pedramktb/schwarzit-probearbeit/migration/fx"
	"github.com/pedramktb/schwarzit-probearbeit/pkg/postgres"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		postgres.FXPostgresModule,
		migrationDI.FXMigrationModule,
		fx.Invoke(func(m migration.Migrator) {
			m.Migrate(context.Background())
		}),
	)
}
