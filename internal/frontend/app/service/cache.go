package service

import (
	"context"

	"github.com/ispringtech/brewkit/internal/backend/api"
)

type ClearCacheParam struct {
	All bool
}

type Cache interface {
	ClearCache(ctx context.Context, param ClearCacheParam) error
}

func NewCacheService(cacheAPI api.CacheAPI) Cache {
	return &cacheService{cacheAPI: cacheAPI}
}

type cacheService struct {
	cacheAPI api.CacheAPI
}

func (service *cacheService) ClearCache(ctx context.Context, param ClearCacheParam) error {
	return service.cacheAPI.ClearCache(ctx, api.ClearParams{
		All: param.All,
	})
}
