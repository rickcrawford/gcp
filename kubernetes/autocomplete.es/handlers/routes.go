package handlers

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/go-chi/chi"
	"google.golang.org/api/iterator"

	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/elastic"
	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/pubsub"
)

// GetRoutes returns the routes for this application
func GetRoutes(esClient *elastic.Client, pubSubClient *pubsub.Client, redisPool *redigo.Pool) http.Handler {
	r := chi.NewRouter()

	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()

		// Sets the name for the new bucket.
		bucketName := "typeahead-catalogs"

		// Creates a Bucket instance.
		bucket := client.Bucket(bucketName)
		query := &storage.Query{}
		it := bucket.Objects(ctx, query)
		for {
			attrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return
			}
			fmt.Fprintln(w, attrs.Name)
		}

	}))

	return r
}
