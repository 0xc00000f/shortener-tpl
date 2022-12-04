package shortener

import (
	"context"

	"github.com/google/uuid"
)

type Shortener interface {
	Short(ctx context.Context, userID uuid.UUID, long string) (short string, err error)
	Get(ctx context.Context, short string) (long string, err error)
	GetAll(ctx context.Context, userID uuid.UUID) (result map[string]string, err error)
}
