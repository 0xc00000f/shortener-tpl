package utils

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

var shortUrlMap urlMap

func init() {
	shortUrlMap = newUrlMap()
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func JoinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

type urlMap map[string]string

func newUrlMap() urlMap {
	return make(urlMap)
}

func EncodeURL(baseUrl string) string {
	encodedUrl := RandStringRunes(6)
	for {
		_, ok := shortUrlMap[encodedUrl]
		if ok {
			encodedUrl = RandStringRunes(6)
		} else {
			break
		}
	}

	shortUrlMap[encodedUrl] = baseUrl

	return encodedUrl
}

func DecodeURL(encodedUrl string) string {
	return shortUrlMap[encodedUrl]
}
