package service

import (
	"context"
	"github.com/koind/cacher/api/internal/domain/repository"
)

// Сервис кэша
type CacheService struct {
	cacheRepository repository.CacheRepositoryInterface
}

// Создает новый сервис кэша
func NewCacheService(cr repository.CacheRepositoryInterface) *CacheService {
	return &CacheService{
		cacheRepository: cr,
	}
}

// Обновить запись, если существует, и создает, если нет
func (s *CacheService) Upsert(ctx context.Context, cache repository.Cache) (*repository.Cache, error) {
	return s.cacheRepository.Upsert(ctx, cache)
}

// Возвращет одну запись по ключу
func (s *CacheService) GetOneByKey(ctx context.Context, key string) (*repository.Cache, error) {
	return s.cacheRepository.GetOneByKey(ctx, key)
}

// Возвращет все записи
func (s *CacheService) GetAll(ctx context.Context) ([]*repository.Cache, error) {
	return s.cacheRepository.GetAll(ctx)
}

// Удаляет одну запись по ключу
func (s *CacheService) Delete(ctx context.Context, key string) error {
	return s.cacheRepository.Delete(ctx, key)
}
