package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"

	_ "github.com/jackc/pgx/stdlib"
)

func Ping(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("pgx", sa.DatabaseAddress)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer cancel()

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
