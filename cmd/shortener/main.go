package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/0xc00000f/shortener-tpl/internal/config"
	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"

	"go.uber.org/zap"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	l, _ := zap.NewProduction()
	defer l.Sync()

	cfg, err := config.New(l)
	if err != nil {
		l.Fatal("creating config error", zap.Error(err))
	}

	encoder := encoder.New(
		encoder.SetStorage(cfg.Storage),
		encoder.SetLength(7),
		encoder.SetLogger(l),
	)

	shortener := shortener.New(
		shortener.SetEncoder(encoder),
		shortener.InitBaseURL(cfg.BaseURL),
		shortener.SetLogger(l),
	)

	router := handlers.NewRouter(shortener)

	l.Info("starting server", zap.String("address", cfg.Address))
	l.Fatal("http server down", zap.Error(http.ListenAndServe(cfg.Address, router)))
}
