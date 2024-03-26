package domain

import "net/url"

type HashKey = string

type URLRepository interface {
	Add(key HashKey, u url.URL)
	GetByHash(key HashKey) (url.URL, bool)
}

type GenShortURLToken = func() HashKey
