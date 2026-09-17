package data

import "github.com/MrAndrey24/cache-proxy/internal/domain"

type cacheRepository struct {
	filePath string
	store    map[string]*domain.CacheResponse
}

func NewCacheRepository(filePath string) *cacheRepository {
	return &cacheRepository{
		filePath: filePath,
		store:    make(map[string]*domain.CacheResponse),
	}
}
