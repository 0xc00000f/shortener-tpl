package handlers

import (
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"

	"github.com/go-chi/chi/v5"
)

func Redirect(sa *shortener.NaiveShortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		short := chi.URLParam(r, "url")

		long, err := sa.Encoder().Get(short)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.Header().Set("Location", long)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
