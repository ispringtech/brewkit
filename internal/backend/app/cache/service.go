package cache

import (
	"context"

	"github.com/ispringtech/brewkit/internal/backend/api"
	"github.com/ispringtech/brewkit/internal/backend/app/docker"
)

func NewCacheService(dockerClient docker.Client) api.CacheAPI {
	return &cacheService{dockerClient: dockerClient}
}

type cacheService struct {
	dockerClient docker.Client
}

func (service *cacheService) ClearCache(ctx context.Context, params api.ClearParams) error {
	return service.dockerClient.ClearCache(ctx, docker.ClearCacheParams{
		All: params.All,
	})
}
