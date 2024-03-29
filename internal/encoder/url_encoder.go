package encoder

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/models"
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

func (ue *URLEncoder) Short(ctx context.Context, userID uuid.UUID, long string) (short string, err error) {
	for {
		short = ue.encode()

		exist, err := ue.storage.IsKeyExist(ctx, short)
		if err != nil {
			return "", err
		}

		if exist {
			continue
		}

		break
	}

	err = ue.storage.Store(ctx, userID, short, long)
	if err != nil {
		var uniqueViolationError *UniqueViolationError

		ok := errors.As(err, &uniqueViolationError)
		if !ok {
			return "", err
		}

		return uniqueViolationError.Short, err
	}

	return short, nil
}

func (ue *URLEncoder) Get(ctx context.Context, short string) (long string, err error) {
	return ue.storage.Get(ctx, short)
}

func (ue *URLEncoder) GetAll(ctx context.Context, userID uuid.UUID) (result map[string]string, err error) {
	return ue.storage.GetAll(ctx, userID)
}

func (ue *URLEncoder) Delete(ctx context.Context, data []models.URL) error {
	return ue.storage.Delete(ctx, data)
}
