package encoder

import (
	"context"

	"github.com/google/uuid"

	"github.com/0xc00000f/shortener-tpl/internal/models"
)

type URLStorager interface {
	Get(ctx context.Context, short string) (string, error)
	GetAll(ctx context.Context, userID uuid.UUID) (result map[string]string, err error)
	Store(ctx context.Context, userID uuid.UUID, short string, long string) error
	IsKeyExist(ctx context.Context, short string) (bool, error)
	Delete(ctx context.Context, data []models.URL) error
}
