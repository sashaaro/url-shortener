package domain

import (
	"errors"
	"net/url"
)

type HashKey = string

type BatchItem struct {
	HashKey HashKey
	URL     url.URL
}

var ErrURLAlreadyExists = errors.New("url already exists")

type URLRepository interface {
	Add(key HashKey, u url.URL) error
	BatchAdd(batch []BatchItem) error
	GetByHash(key HashKey) (url.URL, bool)
}

type GenShortURLToken = func() HashKey
