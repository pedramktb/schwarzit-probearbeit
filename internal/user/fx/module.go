package userDI

import (
	"go.uber.org/fx"

	userDB "github.com/pedramktb/schwarzit-probearbeit/internal/user/db"
)

var FXUserModule = fx.Module("user",
	userDB.FXUserDBProvide,
)
