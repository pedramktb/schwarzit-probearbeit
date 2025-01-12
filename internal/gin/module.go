package ginRouter

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var FXGinRouterModule = fx.Options(
	fx.Provide(
		provideRouter,
		func(r *gin.Engine) gin.IRouter { return r },
	),
	fx.Invoke(run),
)
