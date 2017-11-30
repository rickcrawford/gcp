package elastic

import (
	"strconv"
	"testing"

	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/models"
)

func TestClient(t *testing.T) {
	client, err := NewClient([]string{"http://localhost:9200"}, "elastic", "changeme", "testindex", true)
	if err != nil {
		t.Fatal(err)
	}

	product := models.Product{
		SKU:     1,
		Name:    "test product",
		Content: "123412341231351352124512524",
	}

	err = client.Index(&product)
	if err != nil {
		t.Fatal(err)
	}

	products := make([]models.Product, 10)
	for i := 0; i < 10; i++ {
		products[i] = product
		products[i].SKU += i + 1
		products[i].Name = product.Name + " " + strconv.Itoa(i+1)
	}

	err = client.BulkIndex(products)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Search("test", 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, hit := range res.Hits.Hits {
		t.Log("hit:", hit)
	}

	t.Log("---------------------")
	res, err = client.Autocomplete("te", 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, hit := range res.Hits.Hits {
		t.Log("hit:", hit)
	}

	t.Log("---------------------")
	sug, err := client.Suggest("te", 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sug)

	err = client.DeleteIndex()
	if err != nil {
		t.Fatal(err)
	}

}
