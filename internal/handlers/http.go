package handlers

import (
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(sa *api.ShortenerAPI) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/", SaveURL(*sa))
		r.Post("/api/shorten", SaveURLJson(*sa))

		r.Route("/{url}", func(r chi.Router) {
			r.Get("/", Redirect(*sa))
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "400 page not found", http.StatusBadRequest)
			})
		})

	})
	return r
}
