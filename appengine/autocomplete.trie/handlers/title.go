package handlers

import (
	"net/http"

	"github.com/rickcrawford/gcp/appengine/autocomplete.trie/managers"
)

type titleHandler struct {
	titles managers.Searcher
}

func (t titleHandler) search(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	makeResponse(t.titles, "titles", rw, req)
}

func newTitleHandler(titles managers.Searcher) titleHandler {
	return titleHandler{titles}
}
