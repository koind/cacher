package memory

import (
	"context"
	"github.com/koind/cacher/internal/domain/repository"
	"sync"
)

// Создает фиктивный репозиторий кэша
func NewCacheRepository() *CacheRepository {
	return &CacheRepository{
		DB: make(map[string]*repository.Cache),
	}
}

// Фиктивный репозиторий кэша
type CacheRepository struct {
	sync.RWMutex
	DB map[string]*repository.Cache
}

// Обновить запись, если существует, и создает, если нет
func (r *CacheRepository) Upsert(ctx context.Context, cache repository.Cache) (*repository.Cache, error) {
	_, err := r.GetOneByKey(ctx, cache.Key)
	if err != nil && err != repository.CacheNotFountError {
		return nil, err
	}

	r.Lock()
	defer r.Unlock()

	r.DB[cache.Key] = &cache

	return &cache, nil
}

// Возвращет одну запись по ключу
func (r *CacheRepository) GetOneByKey(ctx context.Context, key string) (*repository.Cache, error) {
	r.RLock()
	defer r.RUnlock()

	cache, has := r.DB[key]
	if !has {
		return nil, repository.CacheNotFountError
	}

	return cache, nil
}

// Возвращет все записи
func (r *CacheRepository) GetAll(ctx context.Context) ([]*repository.Cache, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.DB) <= 0 {
		return nil, nil
	}

	list := make([]*repository.Cache, 0, len(r.DB))

	for _, cache := range r.DB {
		list = append(list, cache)
	}

	return list, nil
}

// Удаляет одну запись по ключу
func (r *CacheRepository) Delete(ctx context.Context, key string) error {
	_, err := r.GetOneByKey(ctx, key)
	if err != nil {
		return err
	}

	r.Lock()
	defer r.Unlock()

	delete(r.DB, key)

	return nil
}
