package api

import (
	"context"
)

type BuildParams struct {
	ForcePull bool
}

type ClearParams struct {
	All bool
}

type BuilderAPI interface {
	Build(ctx context.Context, v Vertex, vars []Var, secretsSrc []SecretSrc, params BuildParams) error
}

type CacheAPI interface {
	ClearCache(ctx context.Context, params ClearParams) error
}
