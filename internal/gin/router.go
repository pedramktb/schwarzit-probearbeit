package ginRouter

import (
	"context"
	"log"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	_ "github.com/pedramktb/schwarzit-probearbeit/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

// @title SchwarzIT Probearbeit API
// @version v1
// @description This is the API for SchwarzIT Probearbeit
// @schemes https
// @securityDefinitions.api_key Bearer

func provideRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}

func run(lc fx.Lifecycle, r *gin.Engine) error {
	var cancel context.CancelFunc
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				ctx, cancel = context.WithCancel(ctx)
				if err := r.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					panic(errors.Wrap(err, "server closed unexpectedly"))
				}
			}()
			log.Println("Server started on :8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cancel != nil {
				cancel()
			}
			return nil
		},
	})

	return nil
}
