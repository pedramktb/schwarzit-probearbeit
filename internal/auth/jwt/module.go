package authJWT

import (
	"github.com/pedramktb/go-base-lib/pkg/env"
	"go.uber.org/fx"
)

var FXAuthJWTProvide = fx.Provide(
	func() *JWT { return create(env.GetOrFail[string]("JWT_SECRET")) },
)
