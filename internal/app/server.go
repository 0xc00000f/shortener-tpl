package app

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0xc00000f/shortener-tpl/internal/utils"
)

var s *http.Server

func init() {
	s = &http.Server{
		Addr:         "http://localhost:8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      serveMuxWithHandlers(),
	}
}

func Server() *http.Server {
	return s
}

func serveMuxWithHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", mainHandler())
	return mux
}

func mainHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			b, err := io.ReadAll(r.Body)
			if err != nil {
				badRequest(w, r)
				return
			}

			longURL := string(b)
			if !utils.IsURL(longURL) {
				badRequest(w, r)
				return
			}

			w.Header().Set("content-type", "raw")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(utils.JoinURL(s.Addr, utils.EncodeURL(longURL))))
			return
		case http.MethodGet:
			urlPathComponent := 1
			urlPart := strings.Split(r.URL.Path, "/")[urlPathComponent]
			originalURL := utils.DecodeURL(urlPart)
			if originalURL == "" {
				badRequest(w, r)
				return
			}

			w.Header().Set("content-type", "application/json")
			w.Header().Set("Location", originalURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		default:
			badRequest(w, r)
		}

	},
	)
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "400 page not found", http.StatusBadRequest)
}
