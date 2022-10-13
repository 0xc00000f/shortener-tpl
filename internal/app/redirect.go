package app

import (
	"net/http"

	helpers2 "github.com/0xc00000f/shortener-tpl/internal/app/helpers"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
)

func Redirect(storage storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPart := chi.URLParam(r, "url")

		originalURL, ok := helpers2.DecodeURLFromStorage(urlPart, storage)
		if !ok {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
