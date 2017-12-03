package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/rickcrawford/gcp/appengine/autocomplete/managers"
)

// GetRoutes returns routes for handlers
func GetRoutes() http.Handler {
	r := chi.NewRouter()

	r.Mount("/products", newProductRouter())
	r.Mount("/search", newSearchRouter())

	return r
}

func newProductRouter() http.Handler {
	r := chi.NewRouter()
	h := productHandler{managers.NewProductManager()}

	r.Post("/", h.save)
	r.Post("/batch", h.saveBatch)

	r.Route("/{productID}", func(r chi.Router) {
		r.Use(h.context)

		r.Get("/", h.get)       // GET /{productID}
		r.Put("/", h.update)    // PUT /{productID}
		r.Delete("/", h.delete) // DELETE /{productID}
	})

	return r
}

func newSearchRouter() http.Handler {
	r := chi.NewRouter()

	h := searchHandler{managers.NewProductManager()}
	r.Get("/", h.search)

	return r
}
