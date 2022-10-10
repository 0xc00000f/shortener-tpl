package app

import (
	helpers2 "github.com/0xc00000f/shortener-tpl/internal/app/helpers"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Redirect(storage storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPart := chi.URLParam(r, "url")

		originalURL, ok := helpers2.DecodeURLFromStorage(urlPart, storage)
		if !ok {
			helpers2.BadRequest(w, r)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
