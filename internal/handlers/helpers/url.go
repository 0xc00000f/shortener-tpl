package helpers

import (
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/utils"
)

func encodeURL(len int) (encodedURL string) {
	encodedURL = utils.RandStringRunes(len)
	return
}

func encodeURLWithDefaultSize() string {
	return encodeURL(6)
}

func EncodeAndStoreURL(baseURL string, urlStorage storage.URLStorage) (encodedURL string) {
	encodedURL = encodeURLWithDefaultSize()
	for {
		_, ok := urlStorage.Get(encodedURL)
		if ok {
			encodedURL = encodeURLWithDefaultSize()
		} else {
			break
		}
	}

	urlStorage.Set(encodedURL, baseURL)
	return encodedURL
}

func DecodeURLFromStorage(encodedURL string, urlStorage storage.URLStorage) (baseURL string, ok bool) {
	return urlStorage.Get(encodedURL)
}
