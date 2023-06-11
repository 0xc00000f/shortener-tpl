package shortener

import (
	"context"

	"github.com/google/uuid"

	"github.com/0xc00000f/shortener-tpl/internal/models"
)

type Shortener interface {
	Short(ctx context.Context, userID uuid.UUID, long string) (short string, err error)
	Get(ctx context.Context, short string) (long string, err error)
	GetAll(ctx context.Context, userID uuid.UUID) (result map[string]string, err error)
	Delete(ctx context.Context, data []models.URL) error
	GetStats(ctx context.Context) (models.Stats, error)
}
