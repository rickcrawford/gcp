package handlers

import (
	"net/http"

	"github.com/rickcrawford/gcp/appengine/autocomplete.trie/managers"
)

type ngramHandler struct {
	ngrams managers.Searcher
}

func (n ngramHandler) search(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	makeResponse(n.ngrams, "ngrams", rw, req)
}

func newNgramHandler(ngrams managers.Searcher) ngramHandler {
	return ngramHandler{ngrams}
}
