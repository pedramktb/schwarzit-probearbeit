package userCache

import (
	"github.com/pedramktb/schwarzit-probearbeit/internal/datasource"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
	"go.uber.org/fx"
)

var FXUserCacheProvide = fx.Provide(
	create,
	fx.Annotate(func(c *cache) datasource.Getter[types.User] { return c }, fx.ResultTags(`name:"cachedUserGetter"`)),
	fx.Annotate(func(c *cache) datasource.Saver[types.User] { return c }, fx.ResultTags(`name:"cachedUserSaver"`)),
	fx.Annotate(func(c *cache) datasource.Deleter[types.User] { return c }, fx.ResultTags(`name:"cachedUserDeleter"`)),
	fx.Annotate(func(c *cache) datasource.UserByEmailGetter { return c }, fx.ResultTags(`name:"cachedUserByEmailGetter"`)),
)
