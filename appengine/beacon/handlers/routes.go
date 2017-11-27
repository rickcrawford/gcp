package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// GetRoutes returns routes for handlers
func GetRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/script.js", scriptHandler)
	r.Get("/info.txt", infoHandler)

	h := beaconHandler{}
	r.Get("/_ah/warmup", h.warmup)
	r.Route("/b.{ext}", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Options("/", h.get)
		r.Head("/", h.get)
		r.Get("/", h.get)
	})
	r.Post("/data", h.post)
	r.Get("/_ah/init", h.warmup)

	return r
}
