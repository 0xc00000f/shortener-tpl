package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"
)

func HealthCheck(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pgxConfig, err := pgxpool.ParseConfig(sa.DatabaseAddress)
		if err != nil {
			sa.L.Error("unable to parsing config", zap.Error(err))
			return
		}

		pgxConnPool, err := pgxpool.ConnectConfig(context.TODO(), pgxConfig)
		if err != nil {
			sa.L.Error("Unable to connect to database", zap.Error(err))
			return
		}
		defer pgxConnPool.Close()

		const defaultTimeout = 2 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		if err := pgxConnPool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
