package datasource

import (
	"context"

	"github.com/google/uuid"

	"github.com/pedramktb/schwarzit-probearbeit/internal/types"
)

type Getter[T any] interface {
	Get(ctx context.Context, id uuid.UUID) (T, error)
}

type Querier[T any] interface {
	Query(ctx context.Context, params types.QueryParams) ([]T, error)
}

type Saver[T any] interface {
	Save(ctx context.Context, t T) (T, error)
}

type Deleter[T any] interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

type VersionGetter[T any] interface {
	GetVersion(ctx context.Context, versionID uuid.UUID) (T, error)
}

type VersionsGetter[T any] interface {
	GetVersions(ctx context.Context, versionID []uuid.UUID) ([]T, error)
}
