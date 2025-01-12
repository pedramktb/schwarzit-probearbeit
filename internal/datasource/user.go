package datasource

import (
	"context"

	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
)

type UserByEmailGetter interface {
	GetByEmail(ctx context.Context, email string) (types.User, error)
}
