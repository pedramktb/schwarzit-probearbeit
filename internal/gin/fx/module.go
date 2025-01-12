package ginDI

import (
	"go.uber.org/fx"

	authGinRouter "github.com/pedramktb/schwarzit-probearbeit/internal/auth/gin"
	ginRouter "github.com/pedramktb/schwarzit-probearbeit/internal/gin"
	userGinRouter "github.com/pedramktb/schwarzit-probearbeit/internal/user/gin"
)

var FXGinRoutersModule = fx.Module("gin",
	ginRouter.FXGinRouterModule,
	authGinRouter.FXAuthGinRouterModule,
	userGinRouter.FXUserGinRouterModule,
)
