package data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/MrAndrey24/cache-proxy/internal/domain"
)

const cacheFile = "cache.json"

type FileCacheRepository struct {
	cache map[string]*domain.CacheResponse
}

func NewFileCacheRepository() *FileCacheRepository {
	return &FileCacheRepository{
		cache: make(map[string]*domain.CacheResponse),
	}
}

func (r *FileCacheRepository) Get(key string) (*domain.CacheResponse, bool) {
	response, hit := r.cache[key]

	return response, hit
}

func (r *FileCacheRepository) Set(key string, response *domain.CacheResponse) {
	r.cache[key] = response
	r.saveCacheFile()
}

func (r *FileCacheRepository) Load() {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return
	}

	loaded := make(map[string]*domain.CacheResponse)
	if err := json.Unmarshal(data, &loaded); err != nil {
		return
	}
	r.cache = loaded
}

func (r *FileCacheRepository) saveCacheFile() {
	data, err := json.MarshalIndent(r.cache, "", " ")
	if err != nil {
		log.Println("warning: could not serialize cache:", err)
		return
	}

	if err := os.WriteFile(cacheFile, data, 0644); err != nil {
		log.Println("warning: could not save cache file:", err)
	}
}

func (r *FileCacheRepository) Clear() error {
	err := os.Remove(cacheFile)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Cache is already empty")
			return nil
		}
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	fmt.Println("Cache cleared successfully")
	r.cache = make(map[string]*domain.CacheResponse)
	return nil
}
