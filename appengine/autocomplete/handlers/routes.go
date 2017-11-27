package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/rickcrawford/gcp/appengine/autocomplete/managers"
	"github.com/rickcrawford/gcp/appengine/autocomplete/middleware"
)

// GetRoutes returns routes for handlers
func GetRoutes() http.Handler {
	r := chi.NewRouter()

	r.Mount("/catalog", newCatalogRouter())
	r.Mount("/search", newSearchRouter())

	return r
}

// NewCatalogRouter creates a completely separate router for administrator routes
func newCatalogRouter() http.Handler {
	h := catalogHandler{managers.NewCatalogManager()}

	r := chi.NewRouter()
	r.Use(middleware.JSONContentType)

	r.Post("/", h.save)

	r.Route("/{catalogID}", func(r chi.Router) {
		r.Use(h.context)
		r.Get("/", h.get) // GET /{catalogID}
		r.Get("/batch", h.getBatch)
		r.Put("/", h.update)    // POST /{catalogID}
		r.Delete("/", h.delete) // DELETE /{catalogID}

		r.Mount("/product", newProductRouter())
		r.Get("/script.js", scriptHandler)
	})

	return r
}

func newProductRouter() http.Handler {
	r := chi.NewRouter()
	h := productHandler{managers.NewProductManager()}

	r.Post("/", h.save)
	r.Post("/batch", h.saveBatch)

	r.Route("/{productID}", func(r chi.Router) {
		r.Use(h.context)
		r.Get("/", h.get)       // GET /{catalogID}/{productID}
		r.Put("/", h.update)    // PUT /{catalogID}/{productID}
		r.Delete("/", h.delete) // DELETE /{catalogID}/{productID}
	})

	return r
}

func newSearchRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/{catalogID}", func(r chi.Router) {
		r.Use(middleware.JSONContentType)

		h := searchHandler{managers.NewProductManager()}
		ch := catalogHandler{managers.NewCatalogManager()}

		// fromage!
		r.Use(ch.context)

		r.Mount("/", http.HandlerFunc(h.search))
	})

	return r
}
