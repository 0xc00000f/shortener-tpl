package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/api"

	"github.com/0xc00000f/shortener-tpl/internal/utils"
)

func SaveURL(s api.Shortener) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		longURL := string(b)
		if err != nil || !utils.IsURL(longURL) {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		short, err := s.Short(longURL)
		if err != nil {
			http.Error(w, "400 page not found", http.StatusBadRequest)
			return
		}

		fullEncodedURL := fmt.Sprintf("http://%s/%s", r.Host, short)
		w.Write([]byte(fullEncodedURL))
	}
}
