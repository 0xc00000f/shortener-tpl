package utils

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

var shortURLMap urlMap

func init() {
	shortURLMap = newURLMap()
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", base, strings.TrimLeft(p, "/"))
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
