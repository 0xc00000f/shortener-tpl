package storage

import (
	"errors"
	"github.com/0xc00000f/shortener-tpl/internal/log"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	ErrNoKeyFound = errors.New("key is not exist")
	ErrEmptyKey   = errors.New("empty string as a key isn't allowed")
	ErrEmptyValue = errors.New("empty string as a value isn't allowed")
)

type MemoryStorage struct {
	storage map[string]string
	history map[uuid.UUID]map[string]string
	l       *zap.Logger
}

func NewMemoryStorage(logger *zap.Logger) MemoryStorage {
	return MemoryStorage{
		storage: make(map[string]string),
		history: make(map[uuid.UUID]map[string]string),
		l:       logger,
	}
}

func (ms MemoryStorage) Get(short string) (long string, err error) {
	ms.l.Info("input", zap.String("short", short))

	if len(short) == 0 {
		err = ErrEmptyKey
		return "", err
	}
	long, ok := ms.storage[short]
	if !ok {
		return "", ErrNoKeyFound
	}

	ms.l.Info("function result", zap.String("long", long), zap.Error(err))
	return long, nil
}

func (ms MemoryStorage) Store(userID uuid.UUID, short, long string) (err error) {
	ms.l.Info("input",
		zap.String("userID", userID.String()),
		zap.String("short", short),
		zap.String("long", long),
		zap.Error(err),
	)

	if len(short) == 0 {
		return ErrEmptyKey
	}
	if len(long) == 0 {
		return ErrEmptyValue
	}
	if userID != uuid.Nil {
		if _, ok := ms.history[userID]; !ok {
			ms.history[userID] = map[string]string{}
		}
		ms.history[userID][short] = long
	}
	ms.storage[short] = long

	ms.l.Info("function result history map", log.MapToFields(ms.history[userID])...)
	ms.l.Info("function result storage map", log.MapToFields(ms.storage)...)
	return nil
}

func (ms MemoryStorage) IsKeyExist(short string) (bool, error) {
	_, ok := ms.storage[short]
	return ok, nil
}

func (ms MemoryStorage) GetAll(uuid uuid.UUID) (result map[string]string, err error) {
	ms.l.Info("function input", zap.String("uuid", uuid.String()))
	result = ms.history[uuid]
	ms.l.Info("function result", log.MapToFields(result)...)
	return result, nil
}
