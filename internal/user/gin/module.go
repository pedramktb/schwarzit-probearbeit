package userGinRouter

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func provideRoutes(e gin.IRouter, r *r, authMiddleware gin.HandlerFunc) {
	g := e.Group("/api/v1/users")
	{
		g.Use(authMiddleware)
		g.POST("/", r.Create)
		g.GET("/", r.Query)
		g.GET("/:id", r.Get)
		g.PUT("/:id", r.Update)
		g.PATCH("/:id", r.Patch)
		g.DELETE("/:id", r.Delete)
		g.GET("/me", r.GetMe)
		g.PUT("/me", r.UpdateMe)
		g.PATCH("/me", r.PatchMe)
		g.DELETE("/me", r.DeleteMe)
	}
}

var FXUserGinRouterModule = fx.Options(
	fx.Provide(fx.Annotate(create, fx.ParamTags(`name:"cachedUserGetter"`, "", `name:"cachedUserSaver"`, `name:"cachedUserDeleter"`))),
	fx.Invoke(fx.Annotate(
		provideRoutes,
		fx.ParamTags("", "", `name:"authMiddleware"`),
	)),
)
