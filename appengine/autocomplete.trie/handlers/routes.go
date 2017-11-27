package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/rickcrawford/gcp/appengine/autocomplete.trie/managers"
	m "github.com/rickcrawford/gcp/appengine/autocomplete.trie/middleware"
)

func GetRoutes(titles, ngrams managers.Searcher) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.DefaultCompress)
	router.Use(middleware.DefaultCompress)
	router.Use(m.JSONContentType)

	router.Get("/titles", newTitleHandler(titles).search)
	router.Get("/ngrams", newNgramHandler(ngrams).search)

	return router
}
