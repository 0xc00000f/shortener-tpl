package shortener

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/workerpool"
)

type NaiveShortener struct {
	encoder       Shortener
	BaseURL       string
	PgxConnPool   *pgxpool.Pool
	Job           chan workerpool.Job
	TrustedSubnet string

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

func SetJobChannel(job chan workerpool.Job) Option {
	return func(ns *NaiveShortener) {
		ns.Job = job
	}
}

func SetTrustedSubnet(subnet string) Option {
	return func(ns *NaiveShortener) {
		ns.TrustedSubnet = subnet
	}
}
