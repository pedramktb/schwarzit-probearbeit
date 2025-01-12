package authGinRouter

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func provideRoutes(e gin.IRouter, r *r) {
	g := e.Group("/auth")
	{
		g.POST("/register", r.Register)
		g.POST("/login", r.Login)
		g.GET("/refresh", r.Refresh)
	}
}

var FXAuthGinRouterModule = fx.Options(
	fx.Provide(
		fx.Annotate(create, fx.ParamTags(`name:"cachedUserByEmailGetter"`, `name:"cachedUserSaver"`, "")),
		fx.Annotate(
			func(r *r) gin.HandlerFunc { return r.AuthMiddleware },
			fx.ResultTags(`name:"authMiddleware"`),
		),
	),
	fx.Invoke(provideRoutes),
)
