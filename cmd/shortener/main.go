package main

import (
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/config"
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/logic"

	"go.uber.org/zap"
)

func main() {
	l, _ := zap.NewProduction()
	defer l.Sync()

	cfg, err := config.NewConfig(l)
	if err != nil {
		l.Fatal("creating config error", zap.Error(err))
	}

	encoder := logic.NewURLEncoder(
		logic.SetStorage(cfg.Storage),
		logic.SetLength(7),
		logic.SetLogger(l),
	)

	shortener := api.NewShortenerAPI(
		api.SetLogic(encoder),
		api.InitBaseURL(cfg.BaseURL),
		api.SetLogger(l),
	)

	router := handlers.NewRouter(shortener)

	l.Info("starting server", zap.String("address", cfg.Address))
	l.Fatal("server fatal error", zap.Error(http.ListenAndServe(cfg.Address, router)))
}
