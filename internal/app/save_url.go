package app

import (
	helpers2 "github.com/0xc00000f/shortener-tpl/internal/app/helpers"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/utils"
	"io"
	"net/http"
)

func SaveURL(storage storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		longURL := string(b)
		if err != nil || !utils.IsURL(longURL) {
			helpers2.BadRequest(w, r)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://" + r.Host + "/" + helpers2.EncodeAndStoreURL(longURL, storage)))
	}
}
