package v1Migration

import (
	"github.com/pedramktb/schwarzit-probearbeit/migration"
	"go.uber.org/fx"
)

var FXV1MigrationProvide = fx.Provide(
	create,
	fx.Annotate(func(m *migrator) migration.Migrator { return m }, fx.ParamTags(`name:"v1Migrator"`)),
)

var FXV1TestMigrationProvide = fx.Provide(
	create,
	fx.Annotate(func(m *migrator) migration.Migrator { return m }),
)
