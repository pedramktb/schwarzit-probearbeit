package userDB

import (
	"github.com/pedramktb/schwarzit-probearbeit/internal/datasource"
	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
	"go.uber.org/fx"
)

var FXUserDBProvide = fx.Provide(
	create,
	func(d *db) datasource.Getter[types.User] { return d },
	func(d *db) datasource.Querier[types.User] { return d },
	func(d *db) datasource.Saver[types.User] { return d },
	func(d *db) datasource.Deleter[types.User] { return d },
	func(d *db) datasource.VersionGetter[types.User] { return d },
	func(d *db) datasource.UserByEmailGetter { return d },
)
