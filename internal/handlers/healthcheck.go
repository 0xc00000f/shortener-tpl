package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"
)

func HealthCheck(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const defaultTimeout = 2 * time.Second

		ctx, cancel := context.WithTimeout(r.Context(), defaultTimeout)
		defer cancel()

		if err := sa.PgxConnPool.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
