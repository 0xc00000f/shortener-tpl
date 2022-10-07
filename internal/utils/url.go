package utils

import (
	"net/url"
)

var shortURLMap urlMap

func init() {
	shortURLMap = newURLMap()
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

type urlMap map[string]string

func newURLMap() urlMap {
	return make(urlMap)
}

func EncodeURL(baseURL string) string {
	encodedURL := RandStringRunes(6)
	for {
		_, ok := shortURLMap[encodedURL]
		if ok {
			encodedURL = RandStringRunes(6)
		} else {
			break
		}
	}

	shortURLMap[encodedURL] = baseURL

	return encodedURL
}

func DecodeURL(encodedURL string) string {
	return shortURLMap[encodedURL]
}
