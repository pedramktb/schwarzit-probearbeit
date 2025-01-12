package main

import (
	"context"

	cmd "github.com/matthiasbruns/golang-base-lib/bootstrap"
	"github.com/pedramktb/schwarzit-probearbeit/migration"
	migrationDI "github.com/pedramktb/schwarzit-probearbeit/migration/fx"
	"github.com/pedramktb/schwarzit-probearbeit/pkg/postgres"
	"go.uber.org/fx"
)

func main() {
	cmd.LoadEnv(".env")
	fx.New(
		postgres.FXPostgresModule,
		migrationDI.FXMigrationModule,
		fx.Invoke(func(m migration.Migrator) {
			m.Migrate(context.Background())
		}),
	)
}
