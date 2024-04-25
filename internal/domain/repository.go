package domain

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/url"
)

type HashKey = string

type BatchItem struct {
	HashKey HashKey
	URL     url.URL
}

type URLEntry struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ErrURLAlreadyExists struct {
	HashKey HashKey
}

func (e *ErrURLAlreadyExists) Error() string {
	return fmt.Sprintf("url %s already exists", e.HashKey)
}

var _ error = (*ErrURLAlreadyExists)(nil)

type URLRepository interface {
	Add(ctx context.Context, key HashKey, u url.URL, userID uuid.UUID) error
	BatchAdd(ctx context.Context, batch []BatchItem, userID uuid.UUID) error
	GetByHash(ctx context.Context, key HashKey) (*url.URL, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]URLEntry, error)
}

type GenShortURLToken = func() HashKey
