package migrationDI

import (
	"github.com/pedramktb/schwarzit-probearbeit/migration"
	v1Migration "github.com/pedramktb/schwarzit-probearbeit/migration/v1"
	"go.uber.org/fx"
)

var FXMigrationModule = fx.Module("migration",
	v1Migration.FXV1MigrationProvide,
	fx.Provide(fx.Annotate(
		func(
			v1Migrator migration.Migrator,
		) migration.Migrator {
			return create(
				v1Migrator,
			)
		},
		fx.ParamTags(`name:"v1Migrator"`),
	)),
)
