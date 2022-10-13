package app

import (
	"fmt"
	"io"
	"net/http"

	helpers2 "github.com/0xc00000f/shortener-tpl/internal/app/helpers"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/utils"
)

func SaveURL(storage storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		longURL := string(b)
		if err != nil || !utils.IsURL(longURL) {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		encodedURL := helpers2.EncodeAndStoreURL(longURL, storage)
		fmt.Fprint(w, "http://%v/%v", r.Host, encodedURL)
	}
}
