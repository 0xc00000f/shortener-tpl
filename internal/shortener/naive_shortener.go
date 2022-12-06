package shortener

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type NaiveShortener struct {
	encoder     Shortener
	BaseURL     string
	PgxConnPool *pgxpool.Pool

	L *zap.Logger
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
		ns.BaseURL = baseURL
	}
}

func SetPgxConnPool(pool *pgxpool.Pool) Option {
	return func(ns *NaiveShortener) {
		ns.PgxConnPool = pool
	}
}

func SetLogger(l *zap.Logger) Option {
	return func(ns *NaiveShortener) {
		ns.L = l
	}
}
