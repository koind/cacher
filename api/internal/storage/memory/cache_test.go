package memory

import (
	"context"
	"github.com/koind/cacher/api/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

const cacheKey = "test-key"

var cacheRepository repository.CacheRepositoryInterface

func init() {
	cacheRepository = NewCacheRepository()
}

func before() {
	cache := repository.Cache{
		Key:   cacheKey,
		Value: "test-value",
	}

	cacheRepository.Upsert(context.Background(), cache)
}

func after() {
	cacheRepository.Delete(context.Background(), cacheKey)
}

func TestCacheRepository_Create(t *testing.T) {
	newCache := repository.Cache{
		Key:   cacheKey,
		Value: "test-value2",
	}

	_, err := cacheRepository.Upsert(context.Background(), newCache)
	assert.Nil(t, err, "не должно быть ошибки при создании")

	cache, _ := cacheRepository.GetOneByKey(context.Background(), cacheKey)
	assert.EqualValues(t, &newCache, cache)

	after()
}

func TestCacheRepository_Update(t *testing.T) {
	before()

	newCache := repository.Cache{
		Key:   cacheKey,
		Value: "test-value3",
	}

	_, err := cacheRepository.Upsert(context.Background(), newCache)
	assert.Nil(t, err, "не должно быть ошибки при обновлении")

	cache, _ := cacheRepository.GetOneByKey(context.Background(), cacheKey)
	assert.EqualValues(t, &newCache, cache)

	after()
}

func TestCacheRepository_GetOneByKey(t *testing.T) {
	before()

	_, err := cacheRepository.GetOneByKey(context.Background(), cacheKey)
	assert.Nil(t, err, "не должно быть ошибки при получении записи")

	after()

	_, err = cacheRepository.GetOneByKey(context.Background(), "test-random-key")
	if assert.NotNil(t, err) {
		assert.EqualError(t, err, repository.CacheNotFountError.Error(), "ошибки должны совподать")
	}
}

func TestCacheRepository_FindAll(t *testing.T) {
	before()

	list, err := cacheRepository.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	after()

	list, _ = cacheRepository.GetAll(context.Background())
	assert.Len(t, list, 0)
}

func TestCacheRepository_Delete(t *testing.T) {
	before()

	err := cacheRepository.Delete(context.Background(), cacheKey)
	assert.Nil(t, err, "не должно быть ошибки при удалении")

	_, err = cacheRepository.GetOneByKey(context.Background(), cacheKey)
	if assert.NotNil(t, err) {
		assert.EqualError(t, err, repository.CacheNotFountError.Error(), "ошибки должны совподать")
	}
}
