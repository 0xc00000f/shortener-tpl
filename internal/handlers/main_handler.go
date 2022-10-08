package handlers

import (
	"github.com/0xc00000f/shortener-tpl/internal/handlers/helpers"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/0xc00000f/shortener-tpl/internal/utils"
	"io"
	"log"
	"net/http"
	"strings"
)

func MainHandler() http.Handler {

	var storage = storage.NewStorage()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			urlPathComponent := 1
			if strings.Split(r.URL.Path, "/")[urlPathComponent] != "" {
				helpers.BadRequest(w, r)
				return
			}
			log.Print(strings.Split(r.URL.Path, "/"))

			b, err := io.ReadAll(r.Body)
			if err != nil {
				helpers.BadRequest(w, r)
				return
			}

			longURL := string(b)
			if !utils.IsURL(longURL) {
				helpers.BadRequest(w, r)
				return
			}

			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("http://" + r.Host + "/" + helpers.EncodeAndStoreURL(longURL, storage)))
			return
		case http.MethodGet:
			urlPathComponent := 1
			urlPart := strings.Split(r.URL.Path, "/")[urlPathComponent]

			originalURL, ok := helpers.DecodeURLFromStorage(urlPart, storage)
			if !ok {
				helpers.BadRequest(w, r)
				return
			}

			w.Header().Set("content-type", "text/plain; charset=utf-8")
			w.Header().Set("Location", originalURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		default:
			helpers.BadRequest(w, r)
		}
	})
}
