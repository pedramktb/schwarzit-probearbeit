package userDI

import (
	"go.uber.org/fx"

	userCache "github.com/pedramktb/schwarzit-probearbeit/internal/user/cache"
	userDB "github.com/pedramktb/schwarzit-probearbeit/internal/user/db"
)

var FXUserModule = fx.Module("user",
	userDB.FXUserDBProvide,
	userCache.FXUserCacheProvide,
)
