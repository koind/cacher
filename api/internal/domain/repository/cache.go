package repository

import (
	"context"
	"errors"
)

var (
	CacheNotFountError = errors.New("не удалось найти запись")
)

// Интерфейс репозитория кэша
type CacheRepositoryInterface interface {
	// Обновить запись, если существует, и создает, если нет
	Upsert(ctx context.Context, cache Cache) (*Cache, error)

	// Возвращет одну запись по ключу
	GetOneByKey(ctx context.Context, key string) (*Cache, error)

	// Возвращет все записи
	GetAll(ctx context.Context) ([]*Cache, error)

	// Удаляет одну запись по ключу
	Delete(ctx context.Context, key string) error
}

// Создает новую структуру кэша
func NewCache(key, value string) *Cache {
	return &Cache{
		Kye:   key,
		Value: value,
	}
}

// Сущность кэша
type Cache struct {
	// Ключ кэша
	Kye string `json:"kye"`
	// Значение кэша
	Value string `json:"value"`
}
