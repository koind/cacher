package service

import (
	"context"
	"github.com/koind/cacher/internal/domain/repository"
	"github.com/koind/cacher/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCacheService_Create(t *testing.T) {
	cacheService := NewCacheService(memory.NewCacheRepository())
	newCache := repository.Cache{Key: "test-key", Value: "test-value"}

	cache, err := cacheService.Upsert(context.Background(), newCache)
	assert.NoError(t, err)
	assert.EqualValues(t, &newCache, cache)
}

func TestCacheService_Update(t *testing.T) {
	cacheService := NewCacheService(memory.NewCacheRepository())
	newCache := repository.Cache{Key: "test-key", Value: "test-value"}

	cache, err := cacheService.Upsert(context.Background(), newCache)
	assert.NoError(t, err)
	assert.EqualValues(t, &newCache, cache)

	changedCache := repository.Cache{Key: "test-key", Value: "test-value2"}
	cache, err = cacheService.Upsert(context.Background(), changedCache)
	assert.NoError(t, err)
	assert.EqualValues(t, &changedCache, cache)
}

func TestCacheService_GetOneByKey(t *testing.T) {
	cacheService := NewCacheService(memory.NewCacheRepository())
	newCache := repository.Cache{Key: "test-key", Value: "test-value"}

	_, err := cacheService.Upsert(context.Background(), newCache)
	assert.NoError(t, err)

	cache, err := cacheService.GetOneByKey(context.Background(), "test-key")
	assert.NoError(t, err)
	assert.EqualValues(t, &newCache, cache)

	_, err = cacheService.GetOneByKey(context.Background(), "test-key1")
	assert.Equal(t, repository.CacheNotFountError, err)
}

func TestCacheService_GetAll(t *testing.T) {
	cacheRepository := memory.NewCacheRepository()
	cacheService := NewCacheService(cacheRepository)
	list, err := cacheService.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, list)

	firstCache := repository.Cache{Key: "test-key1", Value: "test-value1"}
	secondCache := repository.Cache{Key: "test-key2", Value: "test-value2"}
	expectedList := []*repository.Cache{&firstCache, &secondCache}

	_, err = cacheService.Upsert(context.Background(), firstCache)
	assert.NoError(t, err)

	_, err = cacheService.Upsert(context.Background(), secondCache)
	assert.NoError(t, err)

	list, err = cacheService.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, len(expectedList), len(list))

	for _, item := range list {
		if item.Key == firstCache.Key {
			assert.EqualValues(t, &firstCache, item)
		} else {
			assert.EqualValues(t, &secondCache, item)
		}
	}
}

func TestCacheService_Delete(t *testing.T) {
	cacheService := NewCacheService(memory.NewCacheRepository())
	newCache := repository.Cache{Key: "test-key", Value: "test-value"}

	_, err := cacheService.Upsert(context.Background(), newCache)
	assert.NoError(t, err)

	err = cacheService.Delete(context.Background(), "test-key")
	assert.NoError(t, err)
}
