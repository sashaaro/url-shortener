package domain

import (
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

func (e ErrURLAlreadyExists) Error() string {
	return fmt.Sprintf("url %s already exists", e.HashKey)
}

var _ error = (*ErrURLAlreadyExists)(nil)

type URLRepository interface {
	Add(key HashKey, u url.URL) error
	BatchAdd(batch []BatchItem) error
	GetByHash(key HashKey) (url.URL, bool)
}

type GenShortURLToken = func() HashKey
