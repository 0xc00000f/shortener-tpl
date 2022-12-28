package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/0xc00000f/shortener-tpl/internal/config"
	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/rand"
	"github.com/0xc00000f/shortener-tpl/internal/shortener"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/workerpool"

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

	cfg, err := config.New()
	if err != nil {
		l.Fatal("creating config error", zap.Error(err))
	}

	ctx := context.Background()
	pgxConnPool := getPgxConnPool(ctx, cfg.DatabaseAddress)

	urlStorage, err := storage.New(ctx, cfg, pgxConnPool, l)
	if err != nil {
		l.Fatal("creating storage error", zap.Error(err))
	}

	concurrency := 10
	jobsCh := make(chan workerpool.Job, concurrency)
	go func() {
		err := workerpool.RunPoolV2(context.Background(), concurrency, jobsCh)
		if err != nil {
			log.Printf("runpool err: %w", err)
		}
	}()

	urlEncoder := encoder.New(
		encoder.SetStorage(urlStorage),
		encoder.SetLength(ShortLength),
		encoder.SetRandom(rand.New(false)),
		encoder.SetLogger(l),
	)

	urlShortener := shortener.New(
		shortener.SetEncoder(urlEncoder),
		shortener.InitBaseURL(cfg.BaseURL),
		shortener.SetPgxConnPool(pgxConnPool),
		shortener.SetLogger(l),
		shortener.SetJobChannel(jobsCh),
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

func getPgxConnPool(ctx context.Context, connString string) *pgxpool.Pool {
	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil
	}

	pgxConnPool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil
	}

	return pgxConnPool
}
