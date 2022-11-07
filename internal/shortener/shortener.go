package shortener

import (
	"github.com/google/uuid"
)

type Shortener interface {
	Short(userID uuid.UUID, long string) (short string, err error)
	Get(short string) (long string, err error)
	GetAll(userID uuid.UUID) (result map[string]string, err error)
}
