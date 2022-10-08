package helpers

import (
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/utils"
)

func EncodeURL(baseURL string) string {
	encodedURL := utils.RandStringRunes(6)
	for {
		_, ok := storage.Storage.Get(encodedURL)
		if ok {
			encodedURL = utils.RandStringRunes(6)
		} else {
			break
		}
	}

	storage.Storage.Set(encodedURL, baseURL)
	return encodedURL
}

func DecodeURL(encodedURL string) (baseURL string, ok bool) {
	return storage.Storage.Get(encodedURL)
}
