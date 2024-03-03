package domain

import "net/url"

type URLKey = string

type URLRepository interface {
	Add(key URLKey, u url.URL)
	GetByHash(key URLKey) (url.URL, bool)
}
