package handlers

import (
	"compress/flate"
	"net/http"

	"go.uber.org/zap"

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
	r.Use(UnzipBody)

	r.Mount("/debug", middleware.Profiler())

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("400 page not found")); err != nil {
			sa.L.Error("writing body failure", zap.Error(err))
		}
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("400 page not found")); err != nil {
			sa.L.Error("writing body failure", zap.Error(err))
		}
	})

	r.Route("/", func(r chi.Router) {
		r.Post("/", SaveURL(sa))
		r.Get("/ping", HealthCheck(sa))

		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", SaveURLJson(sa))
				r.Post("/batch", Batch(sa))
			})

			r.Route("/user", func(r chi.Router) {
				r.Get("/urls", GetSavedData(sa))
				r.Delete("/urls", Delete(sa))
			})

			if sa.TrustedSubnet != "" {
				r.Get("/internal/stats", GetStats(sa))
			}
		})

		r.Route("/{url}", func(r chi.Router) {
			r.Get("/", Redirect(sa))
		})
	})

	return r
}
