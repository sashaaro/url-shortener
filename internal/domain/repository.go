package domain

import (
	"context"
	"fmt"
	"net/url"
)

type HashKey = string

type BatchItem struct {
	HashKey HashKey
	URL     url.URL
}

type ErrURLAlreadyExists struct {
	HashKey HashKey
}

func (e *ErrURLAlreadyExists) Error() string {
	return fmt.Sprintf("url %s already exists", e.HashKey)
}

var _ error = (*ErrURLAlreadyExists)(nil)

type URLRepository interface {
	Add(ctx context.Context, key HashKey, u url.URL) error
	BatchAdd(ctx context.Context, batch []BatchItem) error
	GetByHash(ctx context.Context, key HashKey) (*url.URL, error)
}

type GenShortURLToken = func() HashKey
