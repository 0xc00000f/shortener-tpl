package main

import (
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/config"
	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/rand"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	"github.com/0xc00000f/shortener-tpl/internal/storage"

	"go.uber.org/zap"
)

func main() {
	l, _ := zap.NewProduction()
	defer l.Sync()

	cfg, err := config.New(l)
	if err != nil {
		l.Fatal("creating config error", zap.Error(err))
	}

	storage, err := storage.New(cfg)
	if err != nil {
		l.Fatal("creating storage error", zap.Error(err))
	}

	encoder := encoder.New(
		encoder.SetStorage(storage),
		encoder.SetLength(7),
		encoder.SetRandom(rand.New(false)),
		encoder.SetLogger(l),
	)

	shortener := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(cfg.BaseURL),
		shortener.SetDatabaseAddress(cfg.DatabaseAddress),
		shortener.SetLogger(l),
	)

	router := handlers.NewRouter(shortener)

	l.Info("starting server", zap.String("address", cfg.Address))
	l.Fatal("http server down", zap.Error(http.ListenAndServe(cfg.Address, router)))
}
