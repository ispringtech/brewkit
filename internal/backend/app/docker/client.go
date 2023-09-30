package docker

import (
	"context"

	"github.com/ispringtech/brewkit/internal/common/maybe"
	"github.com/ispringtech/brewkit/internal/dockerfile"
)

type BuildParams struct {
	Target   string
	SSHAgent maybe.Maybe[string]
	Secrets  []SecretData
	Output   maybe.Maybe[string]
}

type ValueParams struct {
	Var      string
	SSHAgent maybe.Maybe[string]
	Secrets  []SecretData
	UseCache bool
}

type ClearCacheParams struct {
	All bool
}

type SecretData struct {
	ID   string
	Path string
}

type Image struct {
	Img string // Image with repository and tag
}

type Client interface {
	Build(ctx context.Context, dockerfile dockerfile.Dockerfile, params BuildParams) error
	Value(ctx context.Context, dockerfile dockerfile.Dockerfile, params ValueParams) ([]byte, error)
	PullImage(ctx context.Context, img string) error
	ListImages(ctx context.Context, images []string) ([]Image, error)
	BuildImage(ctx context.Context, dockerfilePath string) error

	ClearCache(ctx context.Context, params ClearCacheParams) error
}
