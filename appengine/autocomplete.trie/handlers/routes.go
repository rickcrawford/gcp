package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/rickcrawford/gcp/appengine/autocomplete.trie/managers"
)

func GetRoutes(titles managers.Searcher) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.DefaultCompress)

	router.Get("/", indexHandler)
	router.Get("/script.js", scriptHandler)
	router.Get("/titles", newTitleHandler(titles).search)

	return router
}
