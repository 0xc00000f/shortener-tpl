package encoder

import (
	"github.com/google/uuid"
)

type URLStorager interface {
	Get(short string) (string, error)
	GetAll(userID uuid.UUID) (result map[string]string, err error)
	Store(userID uuid.UUID, short string, long string) error
	IsKeyExist(short string) (bool, error)
}
