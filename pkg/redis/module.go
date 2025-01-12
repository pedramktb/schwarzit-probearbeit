package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/pedramktb/go-base-lib/pkg/env"
	"go.uber.org/fx"
)

var FXRedisModule = fx.Options(fx.Provide(
	func() *redis.Client { return create(env.GetOrFail[string]("REDIS_URI")) },
))
