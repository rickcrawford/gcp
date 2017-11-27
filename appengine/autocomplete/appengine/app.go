// +build !appengine
package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/appengine"

	"github.com/rickcrawford/gcp/appengine/autocomplete/handlers"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.DefaultCompress)

	router.Mount("/", handlers.GetRoutes())

	router.Get("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	// set the application namespaace, and appengine context
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Chi creates a copy of the request, so you need to register the context immediately or
		// you will get an out of flight request context panic
		ctx := appengine.NewContext(r)

		// set the namespace. this should be specific to the logged in user context...
		// for example we would probably want to set this based on the tennat
		// ctx, _ = appengine.Namespace(ctx, "namespace")

		router.ServeHTTP(w, r.WithContext(ctx))
	}))

	appengine.Main()
}
