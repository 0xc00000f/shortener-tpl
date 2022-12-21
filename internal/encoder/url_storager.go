package encoder

import (
	"context"

	"github.com/google/uuid"
)

type URLStorager interface {
	Get(ctx context.Context, short string) (string, error)
	GetAll(ctx context.Context, userID uuid.UUID) (result map[string]string, err error)
	Store(ctx context.Context, userID uuid.UUID, short string, long string) error
	IsKeyExist(ctx context.Context, short string) (bool, error)
}
