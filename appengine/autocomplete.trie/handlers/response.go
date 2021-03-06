package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rickcrawford/gcp/appengine/autocomplete.trie/managers"
)

type Response struct {
	Data     interface{}            `json:"data,omitempty"`
	Errors   []string               `json:"errors,omitempty"`
	Metadata map[string]interface{} `json:"metadata"`
}

func makeResponse(searcher managers.Searcher, name string, rw http.ResponseWriter, req *http.Request) {
	const defaultCount = 5
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	if strings.EqualFold(req.Method, "HEAD") {
		return
	}

	q := strings.TrimSpace(req.FormValue("q"))
	decoder := json.NewEncoder(rw)
	resp := Response{
		Metadata: make(map[string]interface{}),
	}

	resp.Metadata["type"] = name

	if q == "" {
		rw.WriteHeader(http.StatusBadRequest)
		resp.Errors = []string{"no query"}
		decoder.Encode(resp)
		return
	}

	var err error
	count := defaultCount
	countStr := strings.TrimSpace(req.FormValue("count"))
	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			resp.Errors = []string{err.Error()}
			resp.Metadata["status"] = http.StatusInternalServerError
			decoder.Encode(resp)
			return
		}
		if count > 20 || count == 0 {
			count = defaultCount
		}
	}

	resp.Metadata["status"] = http.StatusOK
	resp.Metadata["query"] = q
	resp.Metadata["count"] = count

	results, err := searcher.Search(q, count)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		resp.Errors = []string{err.Error()}
		resp.Metadata["status"] = http.StatusInternalServerError
		decoder.Encode(resp)
		return
	}
	rw.Header().Add("Cache-Control", "max-age=86400")

	resp.Data = results
	decoder.Encode(resp)
}
