package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Category represents a category
type Category struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// Product represents a product
type Product struct {
	SKU          int        `json:"sku"`
	Name         string     `json:"name"`
	Price        float64    `json:"price"`
	UPC          string     `json:"upc"`
	Category     []Category `json:"category"`
	Description  string     `json:"description"`
	Manufacturer string     `json:"manufacturer"`
	Model        string     `json:"model"`
	URL          string     `json:"url"`
	Image        string     `json:"image"`
	Content      string     `json:"content"`
	Updated      time.Time  `json:"updated"`
}

func main() {
	ctx := context.Background()

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
			log.Fatal(err)
		}
		log.Println(attrs.Name)
	}

	obj := bucket.Object("bestbuy/products.json.gz").ReadCompressed(true)
	rdr, err := obj.NewReader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rdr.Close()

	gzr, err := gzip.NewReader(rdr)
	if err != nil {
		log.Fatal(err)
	}
	defer gzr.Close()

	decoder := json.NewDecoder(gzr)
	decoder.UseNumber()
	var v []Product
	err = decoder.Decode(&v)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range v {
		log.Printf("%#v\n", p)
	}

}
