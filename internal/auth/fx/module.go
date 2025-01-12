package authDI

import (
	"go.uber.org/fx"

	authJWT "github.com/pedramktb/schwarzit-probearbeit/internal/auth/jwt"
)

var FXAuthModule = fx.Module("auth",
	authJWT.FXAuthJWTProvide,
)
