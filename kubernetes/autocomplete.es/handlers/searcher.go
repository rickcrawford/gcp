package handlers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	redigo "github.com/garyburd/redigo/redis"
	es "gopkg.in/olivere/elastic.v5"

	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/elastic"
	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/models"
)

const defaultExpiresSeconds = 300
const defaultCount = 5

func DigestString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

type searcher struct {
	esClient *elastic.Client
	pool     *redigo.Pool
}

func (s searcher) search(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	if strings.EqualFold(req.Method, "HEAD") {
		return
	}

	searchType, _ := strconv.Atoi(req.FormValue("type"))
	switch searchType {
	case 1:
		writeResult(rw, req, s.pool, "suggest", s.esClient.Suggest)

	case 2:
		writeResult(rw, req, s.pool, "autocomplete", s.esClient.Autocomplete)

	case 3:
		writeResult(rw, req, s.pool, "search", s.esClient.Search)

	default:
		writeResult(rw, req, s.pool, "prefix", s.esClient.Prefix)

	}

}

func writeResult(rw http.ResponseWriter, req *http.Request, pool *redigo.Pool, typeName string, fn func(string, int) (*es.SearchResult, error)) {
	query := req.FormValue("q")
	count, _ := strconv.Atoi(req.FormValue("c"))
	if count == 0 {
		count = defaultCount
	}

	resp := models.Response{
		Metadata: map[string]interface{}{
			"count": count,
			"query": query,
			"type":  typeName,
		},
	}

	var etag string

	conn := pool.Get()
	defer conn.Close()

	var err error
	var result interface{}

	key := fmt.Sprintf("%s:%s:%d", typeName, query, count)

	log.Println("key", key)

	if result, err = conn.Do("GET", key); err != nil || result == nil {
		var searchResult *es.SearchResult
		if searchResult, err = fn(query, count); err == nil {
			var products []models.Product

			switch typeName {
			case "suggest":
				for _, value := range searchResult.Suggest {
					for _, option := range value {
						products = make([]models.Product, len(option.Options))
						for i, hit := range option.Options {
							json.Unmarshal(*hit.Source, &products[i])
						}
					}
				}

			default:
				products = make([]models.Product, len(searchResult.Hits.Hits))
				for i, hit := range searchResult.Hits.Hits {
					json.Unmarshal(*hit.Source, &products[i])
				}

			}

			var data []byte
			data, err = json.Marshal(products)
			result = data
			_, err = conn.Do("SETEX", key, defaultExpiresSeconds, result)

		}
	} else {
		resp.Metadata["cached"] = true
	}
	if err != nil {
		resp.Errors = []string{err.Error()}
	} else {
		rw.Header().Add("Cache-Control", "max-age=300")
		etag = DigestString(string([]byte(result.([]byte))))
		resp.Metadata["etag"] = etag
		if etag != "" && req.Header.Get("If-None-Match") == etag {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		rw.Header().Add("ETag", etag)

	}

	err = json.Unmarshal([]byte(result.([]byte)), &resp.Data)

	status := 200
	if len(resp.Errors) > 0 {
		status = 500
	}

	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(resp)
}
