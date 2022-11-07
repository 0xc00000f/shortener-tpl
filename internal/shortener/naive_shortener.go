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

type Option func(ns *NaiveShortener)

func New(options ...Option) *NaiveShortener {
	ns := NaiveShortener{}

	for _, fn := range options {
		fn(&ns)
	}

	return &ns
}

func (ns *NaiveShortener) Encoder() Shortener {
	return ns.encoder
}

func SetEncoder(encoder Shortener) Option {
	return func(ns *NaiveShortener) {
		ns.encoder = encoder
	}
}

func InitBaseURL(baseURL string) Option {
	return func(ns *NaiveShortener) {
		if len(baseURL) > 0 {
			ns.BaseURL = baseURL
			return
		}

		ns.BaseURL = os.Getenv("BASE_URL")
	}
}

func SetLogger(l *zap.Logger) Option {
	return func(ns *NaiveShortener) {
		ns.L = l
	}
}
