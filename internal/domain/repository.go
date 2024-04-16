package domain

import "net/url"

type HashKey = string

type BatchItem struct {
	HashKey HashKey
	URL     url.URL
}

type URLRepository interface {
	Add(key HashKey, u url.URL)
	BatchAdd(batch []BatchItem)
	GetByHash(key HashKey) (url.URL, bool)
}

type GenShortURLToken = func() HashKey
