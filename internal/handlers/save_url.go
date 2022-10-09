package handlers

import (
	"github.com/0xc00000f/shortener-tpl/internal/handlers/helpers"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/utils"
	"io"
	"net/http"
)

func SaveURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	longURL := string(b)
	if err != nil || !utils.IsURL(longURL) {
		helpers.BadRequest(w, r)
		return
	}

	w.Header().Set("content-type", "raw")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://" + r.Host + "/" + helpers.EncodeAndStoreURL(longURL, storage.Storage)))
	return
}
