package managers

import (
	"compress/gzip"
	"context"
	"encoding/json"

	"cloud.google.com/go/storage"
	"github.com/fvbock/trie"
)

type Category struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Product struct {
	SKU      int
	Name     string
	Price    float64
	UPC      string
	Category []Category
	// Shipping     float64
	Description  string
	Manufacturer string
	Model        string
	URL          string
	Image        string
}

type productSearcher struct {
	lookup map[string]Product
	trie   *trie.Trie
}

func (p productSearcher) Search(query string, count int) (*Result, error) {
	prefix := FormatProductKey(query, "_")

	keywords := make([]Keyword, 0)
	total := 0
	for _, member := range p.trie.PrefixMembers(prefix) {
		if total == count {
			break
		}
		if product, isPresent := p.lookup[member.Value]; isPresent {
			keywords = append(keywords, Keyword{
				Value: product.Name,
				Count: member.Count,
			})
			total++
		}
	}

	result := &Result{
		Query:    prefix,
		Keywords: keywords,
	}
	return result, nil
}

func ProductSearcher(bucketName, path string) (Searcher, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(path)
	rdr, err := obj.ReadCompressed(true).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	gzrdr, err := gzip.NewReader(rdr)
	if err != nil {
		return nil, err
	}
	defer gzrdr.Close()

	var products []Product
	decoder := json.NewDecoder(gzrdr)
	decoder.UseNumber()

	err = decoder.Decode(&products)
	if err != nil {
		return nil, err
	}

	trie := trie.NewTrie()
	lookup := make(map[string]Product)
	for _, product := range products {
		key := FormatProductKey(product.Name, "_")
		lookup[key] = product
		trie.Add(key)
	}

	searcher := &productSearcher{
		trie:   trie,
		lookup: lookup,
	}

	return searcher, nil
}
