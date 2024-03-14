package adapters

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/sashaaro/url-shortener/internal/domain"
)

func GenBase64ShortURLToken() domain.HashKey {
	length := 8
	bufSize := length*6/8 + 1
	buf := make([]byte, bufSize)
	n, err := rand.Read(buf)
	if err != nil || n != bufSize {
		panic(fmt.Errorf("error while retriving random data: %d %v", n, err.Error()))
	}
	return base64.URLEncoding.EncodeToString(buf)[:length]
}
