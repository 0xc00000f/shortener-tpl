package encoder

import (
	"github.com/0xc00000f/shortener-tpl/internal/rand"

	"go.uber.org/zap"
)

type URLEncoder struct {
	length  int
	storage URLStorager
	rand    rand.Random

	l *zap.Logger
}

type Option func(ue *URLEncoder)

func New(options ...Option) *URLEncoder {
	const preferredLength = 6
	ue := URLEncoder{length: preferredLength}

	for _, fn := range options {
		fn(&ue)
	}

	return &ue
}

func SetLength(length int) Option {
	return func(ue *URLEncoder) {
		ue.length = length
	}
}

func SetStorage(s URLStorager) Option {
	return func(ue *URLEncoder) {
		ue.storage = s
	}
}

func SetLogger(l *zap.Logger) Option {
	return func(ue *URLEncoder) {
		ue.l = l
	}
}

func SetRandom(r rand.Random) Option {
	return func(ue *URLEncoder) {
		ue.rand = r
	}
}

func (ue *URLEncoder) encode() string {
	return ue.rand.String(ue.length)
}

func (ue *URLEncoder) Short(long string) (short string, err error) {
	for {
		short = ue.encode()
		exist, err := ue.storage.IsKeyExist(short)
		if err != nil {
			return "", err
		}
		if exist {
			continue
		}
		break
	}

	err = ue.storage.Store(short, long)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (ue *URLEncoder) Get(short string) (long string, err error) {
	return ue.storage.Get(short)
}
