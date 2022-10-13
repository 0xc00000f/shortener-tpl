package helpers

import (
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/utils"
)

func encodeURL(len int) string {
	return utils.RandStringRunes(len)
}

func encodeURLWithDefaultSize() string {
	return encodeURL(6)
}

func EncodeAndStoreURL(baseURL string, urlStorage storage.URLStorage) (encodedURL string) {
Loop:
	for {
		encodedURL = encodeURLWithDefaultSize()
		_, ok := urlStorage.Get(encodedURL)
		if ok {
			continue Loop
		}
		break
	}

	urlStorage.Set(encodedURL, baseURL)
	return encodedURL
}

func DecodeURLFromStorage(encodedURL string, urlStorage storage.URLStorage) (baseURL string, ok bool) {
	return urlStorage.Get(encodedURL)
}
