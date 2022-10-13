package handlers

import (
	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/logic"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(sa *api.ShortenerApi) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Post("/", SaveURL(sa.Logic()))

		r.Route("/{url}", func(r chi.Router) {
			r.Get("/", Redirect(sa.Logic()))
			r.Post("/", logic.BadRequest)
		})

	})
	return r
}
