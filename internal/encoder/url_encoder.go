package encoder

import (
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/rand"
)

const PreferredLength = 6

type URLEncoder struct {
	length  int
	storage URLStorager
	rand    rand.Random

	l *zap.Logger
}

type Option func(ue *URLEncoder)

func New(options ...Option) *URLEncoder {
	ue := URLEncoder{length: PreferredLength}

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

func (ue *URLEncoder) Short(userID uuid.UUID, long string) (short string, err error) {
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

	err = ue.storage.Store(userID, short, long)
	if err != nil {
		var uniqueViolationError *UniqueViolationError
		if errors.As(err, &uniqueViolationError) {
			if err, ok := err.(*UniqueViolationError); ok { //nolint:errorlint
				return err.Short, err
			}
		}

		return "", err
	}

	return short, nil
}

func (ue *URLEncoder) Get(short string) (long string, err error) {
	return ue.storage.Get(short)
}

func (ue *URLEncoder) GetAll(userID uuid.UUID) (result map[string]string, err error) {
	return ue.storage.GetAll(userID)
}
