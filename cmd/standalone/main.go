package main

import (
	"go.uber.org/fx"

	authDI "github.com/pedramktb/schwarzit-probearbeit/internal/auth/fx"
	ginDI "github.com/pedramktb/schwarzit-probearbeit/internal/gin/fx"
	userDI "github.com/pedramktb/schwarzit-probearbeit/internal/user/fx"
	"github.com/pedramktb/schwarzit-probearbeit/pkg/postgres"
	"github.com/pedramktb/schwarzit-probearbeit/pkg/redis"
)

var app *fx.App

func init() {
	app = fx.New(
		postgres.FXPostgresModule,
		redis.FXRedisModule,
		authDI.FXAuthModule,
		userDI.FXUserModule,
		ginDI.FXGinRoutersModule,
	)
}

func main() {
	app.Run()
}
