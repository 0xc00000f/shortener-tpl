package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/0xc00000f/shortener-tpl/internal/config"
	"github.com/0xc00000f/shortener-tpl/internal/encoder"

	"go.uber.org/zap"
)

func New(
	ctx context.Context,
	cfg config.Cfg,
	pgxConnPool *pgxpool.Pool,
	l *zap.Logger,
) (encoder.URLStorager, error) {
	if pgxConnPool != nil {
		return NewDatabaseStorage(ctx, pgxConnPool, l)
	}

	if len(cfg.Filepath) == 0 {
		return NewMemoryStorage(l), nil
	}

	storage, err := NewFileStorage(cfg.Filepath, l)
	if err != nil {
		l.Error("creating file storage err", zap.Error(err))
		return nil, err
	}

	err = storage.InitMemory()
	if err != nil {
		l.Error("init file storage memory err", zap.Error(err))
		return nil, err
	}

	return storage, nil
}
