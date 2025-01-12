package migration

import "context"

type Migrator interface {
	Migrate(ctx context.Context)
}
