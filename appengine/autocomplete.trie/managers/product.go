package managers

import (
	"compress/gzip"
	"context"
	"encoding/json"

	"cloud.google.com/go/storage"
	"github.com/fvbock/trie"

	"github.com/rickcrawford/gcp/common/models"
)

type productSearcher struct {
	lookup map[string]models.Product
	trie   *trie.Trie
}

func (p productSearcher) Search(query string, count int) ([]models.Product, error) {
	prefix := FormatProductKey(query, "_")

	products := make([]models.Product, 0)
	total := 0
	for _, member := range p.trie.PrefixMembers(prefix) {
		if total == count {
			break
		}
		if product, isPresent := p.lookup[member.Value]; isPresent {
			products = append(products, product)
			total++
		}
	}

	return products, nil
}

func ProductSearcher(bucketName, productPath, imagesPath string) (Searcher, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)
	objProd := bucket.Object(productPath)
	rdr, err := objProd.ReadCompressed(true).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	gzrdr, err := gzip.NewReader(rdr)
	if err != nil {
		return nil, err
	}
	defer gzrdr.Close()

	var products []models.Product
	decoder := json.NewDecoder(gzrdr)
	decoder.UseNumber()

	err = decoder.Decode(&products)
	if err != nil {
		return nil, err
	}

	objImages := bucket.Object(imagesPath)
	rdr, err = objImages.ReadCompressed(true).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	gzrdr, err = gzip.NewReader(rdr)
	if err != nil {
		return nil, err
	}
	defer gzrdr.Close()

	var images []models.Image
	decoder = json.NewDecoder(gzrdr)
	decoder.UseNumber()

	err = decoder.Decode(&images)
	if err != nil {
		return nil, err
	}

	imageSku := make(map[int]models.Image, len(images))
	for _, image := range images {
		imageSku[image.SKU] = image
	}

	trie := trie.NewTrie()
	lookup := make(map[string]models.Product)
	for i := range products {
		if image, isPresent := imageSku[products[i].SKU]; isPresent {
			products[i].Content = image.Content
			products[i].Image = image.Image
		}
		key := FormatProductKey(products[i].Name, "_")
		lookup[key] = products[i]
		trie.Add(key)
	}

	searcher := &productSearcher{
		trie:   trie,
		lookup: lookup,
	}

	return searcher, nil
}
