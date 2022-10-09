package app

import (
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/handlers/helpers"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	storage := storage.NewStorage()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.SaveURL(storage))

		r.Route("/{url}", func(r chi.Router) {
			r.Get("/", handlers.Redirect(storage))
			r.Post("/", helpers.BadRequest)
		})

	})
	return r
}
