package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	redigo "github.com/garyburd/redigo/redis"
	es "gopkg.in/olivere/elastic.v5"

	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/elastic"
)

const defaultExpiresSeconds = 300

// Response is a response struct for results
type Response struct {
	Data     interface{}            `json:"data,omitempty"`
	Errors   []string               `json:"errors,omitempty"`
	Metadata map[string]interface{} `json:"metadata"`
}

type searcher struct {
	esClient *elastic.Client
	pool     *redigo.Pool
}

func (s searcher) search(rw http.ResponseWriter, req *http.Request) {
	writeResult(rw, req, s.pool, "search", s.esClient.Search)
}

func (s searcher) autocomplete(rw http.ResponseWriter, req *http.Request) {
	writeResult(rw, req, s.pool, "autocomplete", s.esClient.Autocomplete)
}

func (s searcher) suggest(rw http.ResponseWriter, req *http.Request) {
	writeResult(rw, req, s.pool, "suggest", s.esClient.Suggest)
}

func writeResult(rw http.ResponseWriter, req *http.Request, pool *redigo.Pool, typeName string, fn func(string, int) (*es.SearchResult, error)) {
	query := req.FormValue("q")
	count, _ := strconv.Atoi(req.FormValue("c"))
	if count == 0 {
		count = 10
	}

	resp := Response{
		Metadata: map[string]interface{}{
			"count": count,
			"query": query,
			"type":  typeName,
		},
	}

	conn := pool.Get()
	defer conn.Close()

	var err error
	var result interface{}

	key := fmt.Sprintf("%s:%s:%d", typeName, query, count)

	log.Println("key", key)

	if result, err = conn.Do("GET", key); err != nil || result == nil {
		var searchResult interface{}
		if searchResult, err = fn(query, count); err == nil {
			log.Println("search result", typeName, query, count, searchResult)

			var data []byte
			data, err = json.Marshal(searchResult)
			result = data
			_, err = conn.Do("SETEX", key, defaultExpiresSeconds, result)
		}
	} else {
		resp.Metadata["cached"] = true
	}
	if err != nil {
		resp.Errors = []string{err.Error()}
	}
	err = json.Unmarshal([]byte(result.([]byte)), &resp.Data)

	status := 200
	if len(resp.Errors) > 0 {
		status = 500
	}

	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(resp)
}
