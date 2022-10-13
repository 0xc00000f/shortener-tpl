package handlers

import (
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/api"

	"github.com/go-chi/chi/v5"
)

func Redirect(s api.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPart := chi.URLParam(r, "url")

		originalURL, err := s.Get(urlPart)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
