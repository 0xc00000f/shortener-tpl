package handlers

import (
	"github.com/0xc00000f/shortener-tpl/internal/handlers/helpers"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	urlPart := chi.URLParam(r, "url")

	originalURL, ok := helpers.DecodeURLFromStorage(urlPart, storage.Storage)
	if !ok {
		helpers.BadRequest(w, r)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return
}
