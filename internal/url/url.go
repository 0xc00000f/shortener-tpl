package url

import (
	"errors"
	"net/url"
)

var ErrInvalidURL = errors.New("invalid url")

func Valid(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}
