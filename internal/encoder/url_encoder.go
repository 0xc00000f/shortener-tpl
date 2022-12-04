package encoder

import (
	"errors"
	"fmt"

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
			return "", fmt.Errorf("storage creating short failure: %w", err)
		}

		if exist {
			continue
		}

		break
	}

	var uniqueViolationError *UniqueViolationError

	err = ue.storage.Store(userID, short, long)
	if errors.As(err, &uniqueViolationError) {
		if err, ok := err.(*UniqueViolationError); ok {
			return err.Short, err
		}
	}
	if err != nil { //nolint:wsl
		return "", fmt.Errorf("storage creating short failure: %w", err)
	}

	return short, nil
}

func (ue *URLEncoder) Get(short string) (long string, err error) {
	long, err = ue.storage.Get(short)
	if err != nil {
		return "", fmt.Errorf("storage get failure: %w", err)
	}

	return long, nil
}

func (ue *URLEncoder) GetAll(userID uuid.UUID) (result map[string]string, err error) {
	result, err = ue.storage.GetAll(userID)
	if err != nil {
		return nil, fmt.Errorf("storage get all failure: %w", err)
	}

	return result, nil
}
