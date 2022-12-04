package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/0xc00000f/shortener-tpl/internal/config"
	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/rand"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	"github.com/0xc00000f/shortener-tpl/internal/storage"

	"go.uber.org/zap"
)

const (
	ShortLength              = 7
	defaultReadHeaderTimeout = 3 * time.Second
)

func main() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer l.Sync() //nolint:errcheck,wsl

	cfg, err := config.New(l)
	if err != nil {
		l.Fatal("creating config error", zap.Error(err))
	}

	urlStorage, err := storage.New(context.Background(), cfg)
	if err != nil {
		l.Fatal("creating storage error", zap.Error(err))
	}

	urlEncoder := encoder.New(
		encoder.SetStorage(urlStorage),
		encoder.SetLength(ShortLength),
		encoder.SetRandom(rand.New(false)),
		encoder.SetLogger(l),
	)

	urlShortener := shortener.New(
		shortener.SetEncoder(urlEncoder),
		shortener.InitBaseURL(cfg.BaseURL),
		shortener.SetDatabaseAddress(cfg.DatabaseAddress),
		shortener.SetLogger(l),
	)

	router := handlers.NewRouter(urlShortener)
	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           router,
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	l.Info("starting server", zap.String("address", cfg.Address))
	l.Fatal("http server down", zap.Error(server.ListenAndServe()))
}
