package postgres

import (
	"github.com/pedramktb/go-base-lib/pkg/env"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var FXPostgresModule = fx.Options(fx.Provide(
	func() *gorm.DB { return create(env.GetOrFail[string]("POSTGRES_URI")) },
))
