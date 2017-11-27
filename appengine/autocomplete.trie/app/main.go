package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rickcrawford/autocomplete.trie/handlers"
	"github.com/rickcrawford/autocomplete.trie/managers"
)

func main() {
	bucketName := os.Getenv("BUCKET_NAME")

	productsPath := os.Getenv("PRODUCTS_PATH")
	products, err := managers.ProductSearcher(bucketName, productsPath)
	if err != nil {
		log.Fatalf("Error products from bucket: %s", err)
	}

	keywordsPath := os.Getenv("KEYWORDS_PATH")
	keywords, err := managers.KeywordSearcher(bucketName, keywordsPath)
	if err != nil {
		log.Fatalf("Error keywords from bucket: %s", err)
	}

	r := handlers.GetRoutes(products, keywords)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
