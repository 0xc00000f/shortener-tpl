package app

import (
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"net/http"
	"time"
)

var s *http.Server

func init() {
	s = &http.Server{
		Addr:         "localhost:8080",
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
	mux.Handle("/", handlers.MainHandler())
	return mux
}
