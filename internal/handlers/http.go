package handlers

import (
	"compress/flate"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/shortener"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(sa *shortener.NaiveShortener) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentEncoding("gzip"))
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)
	r.Use(CookieAuth)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("400 page not found"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("400 page not found"))
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/", SaveURL(sa))
		r.Get("/ping", Ping(sa))

		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", SaveURLJson(sa))
				r.Post("/batch", Batch(sa))
			})

			r.Get("/user/urls", GetSavedData(sa))
		})

		r.Route("/{url}", func(r chi.Router) {
			r.Get("/", Redirect(sa))
		})
	})

	return r
}
