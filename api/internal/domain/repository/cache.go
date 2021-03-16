package repository

import "context"

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

// Сущность кэша
type Cache struct {
	// Ключ кэша
	Kye string `json:"kye"`
	// Значение кэша
	Value string `json:"value"`
}
