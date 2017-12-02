package handlers

import (
	"net/http"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/go-chi/chi"

	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/elastic"
)

// GetRoutes returns the routes for this application
func GetRoutes(esClient *elastic.Client, pool *redigo.Pool) http.Handler {
	r := chi.NewRouter()

	searchHandler := searcher{
		esClient: esClient,
		pool:     pool,
	}

	r.Get("/search", searchHandler.search)
	r.Get("/autocomplete", searchHandler.autocomplete)
	r.Get("/suggest", searchHandler.suggest)

	// r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	ctx := r.Context()

	// 	client, err := storage.NewClient(ctx)
	// 	if err != nil {
	// 		log.Fatalf("Failed to create client: %v", err)
	// 	}
	// 	defer client.Close()

	// 	// Sets the name for the new bucket.
	// 	bucketName := "typeahead-catalogs"

	// 	// Creates a Bucket instance.
	// 	bucket := client.Bucket(bucketName)
	// 	query := &storage.Query{}
	// 	it := bucket.Objects(ctx, query)
	// 	for {
	// 		attrs, err := it.Next()
	// 		if err == iterator.Done {
	// 			break
	// 		}
	// 		if err != nil {
	// 			return
	// 		}
	// 		fmt.Fprintln(w, attrs.Name)
	// 	}
	// }))

	return r
}
