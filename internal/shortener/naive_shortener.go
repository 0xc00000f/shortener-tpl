package shortener

import (
	"os"

	"go.uber.org/zap"
)

type NaiveShortener struct {
	encoder Shortener
	BaseURL string
	L       *zap.Logger
}

type Option func(sa *NaiveShortener)

func New(options ...Option) *NaiveShortener {
	sa := NaiveShortener{}

	for _, fn := range options {
		fn(&sa)
	}

	return &sa
}

func (sa *NaiveShortener) Encoder() Shortener {
	return sa.encoder
}

func SetEncoder(encoder Shortener) Option {
	return func(sa *NaiveShortener) {
		sa.encoder = encoder
	}
}

func InitBaseURL(baseURL string) Option {
	return func(sa *NaiveShortener) {
		if len(baseURL) > 0 {
			sa.BaseURL = baseURL
			return
		}

		sa.BaseURL = os.Getenv("BASE_URL")
	}
}

func SetLogger(l *zap.Logger) Option {
	return func(sa *NaiveShortener) {
		sa.L = l
	}
}
