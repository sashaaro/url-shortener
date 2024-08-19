package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/url"
)

// ключ для короткой ссылки
type HashKey = string

// структура для обновления
type BatchItem struct {
	HashKey HashKey
	URL     url.URL
}

// ссылка короткая, оригинал
type URLEntry struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ошибка ссылка была удалена
var ErrURLDeleted = fmt.Errorf("url deleted")

// ошибка ссылка уже существует
type ErrURLAlreadyExists struct {
	HashKey HashKey
}

// имлементация error
func (e *ErrURLAlreadyExists) Error() string {
	return fmt.Sprintf("url %s already exists", e.HashKey)
}

var _ error = (*ErrURLAlreadyExists)(nil)

// основной интерфейс управления ссылками
type URLRepository interface {
	Add(ctx context.Context, key HashKey, u url.URL, userID uuid.UUID) error
	BatchAdd(ctx context.Context, batch []BatchItem, userID uuid.UUID) error
	GetByHash(ctx context.Context, key HashKey) (*url.URL, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]URLEntry, error)
	DeleteByUser(ctx context.Context, keys []HashKey, userID uuid.UUID) (bool, error)
}

// интерфейс получения
type GenShortURLToken = func() HashKey
