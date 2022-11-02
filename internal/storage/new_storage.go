package storage

import (
	"github.com/0xc00000f/shortener-tpl/internal/config"
	"github.com/0xc00000f/shortener-tpl/internal/encoder"

	"go.uber.org/zap"
)

func New(cfg config.Cfg) (encoder.URLStorager, error) {
	if len(cfg.Filepath) == 0 {
		return NewMemoryStorage(cfg.L), nil
	}

	storage, err := NewFileStorage(cfg.Filepath, cfg.L)
	if err != nil {
		cfg.L.Error("creating file storage err", zap.Error(err))
		return nil, err
	}

	err = storage.InitMemory()
	if err != nil {
		cfg.L.Error("init file storage memory err", zap.Error(err))
		return nil, err
	}
	return storage, nil
}