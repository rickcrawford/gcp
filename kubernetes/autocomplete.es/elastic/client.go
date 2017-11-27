package elastic

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/rickcrawford/autocomplete.es/models"
	elastic "gopkg.in/olivere/elastic.v5"
)

const mappingType = "product"

const mapping = `
{
	"settings":{
		"number_of_shards": 3,
		"number_of_replicas": 1,
		"analysis": {
	        "filter": {
	            "autocomplete_filter": { 
	                "type":     "edge_ngram",
	                "min_gram": 1,
	                "max_gram": 20
	            }
	        },
	        "analyzer": {
	            "autocomplete": {
	                "type":      "custom",
	                "tokenizer": "standard",
	                "filter": [
	                    "lowercase",
	                    "autocomplete_filter" 
	                ]
	            }
	        }
	    }
	},
	"mappings":{
		"product":{
			"properties":{
				"sku":{
					"type":"integer"
				},
				"name":{
					"type":"text",
					"copy_to":"name_autocomplete"
				},
				"price": {
					"type":"float"
				},
				"upc":{
					"type":"text",
					"index":"not_analyzed"
				},
				"category":{
					"type":"object"
				},
				"description":{
					"type":"text",
					"store": true
				},
				"manufacturer":{
					"type":"text"
				},				
				"model":{
					"type":"text",
					"index":"not_analyzed"
				},
				"url":{
					"type":"text",
					"store":true,
					"index":false
				},
				"image":{
					"type":"text",
					"store":true,
					"index":false
				},
				"content":{
					"type":"text",
					"store":true,
					"index":false
				},
				"weight": {
					"type":"float"
				},
				"updated":{
					"type":"date"
				},
				"suggestion":{
					"type":"completion"
				},
				"name_autocomplete":{
					"type":"text",
					 "analyzer": "autocomplete"
				}
			}
		}
	}
}`

// Client is an ES client for searching our application
type Client struct {
	client    *elastic.Client
	indexName string
}

// Index product
func (c *Client) Index(product *models.Product) error {
	ctx := context.Background()
	ID := strconv.Itoa(product.SKU)
	product.Updated = time.Now().UTC()
	if product.Suggestion == nil {
		product.Suggestion = []models.Suggestion{
			{
				Input:  product.Name,
				Weight: 1,
			},
		}
	}

	_, err := c.client.Index().
		Index(c.indexName).
		Type(mappingType).
		Id(ID).
		BodyJson(product).
		Do(ctx)

	if err == nil {
		_, err = c.client.Flush().Index(c.indexName).Do(ctx)
	}
	return err
}

// BulkIndex products
func (c *Client) BulkIndex(products []models.Product) error {
	ctx := context.Background()
	bulkReq := c.client.Bulk().Index(c.indexName)

	for i := range products {
		ID := strconv.Itoa(products[i].SKU)
		products[i].Updated = time.Now().UTC()
		if products[i].Suggestion == nil {
			products[i].Suggestion = []models.Suggestion{
				{
					Input:  products[i].Name,
					Weight: 1,
				},
			}
		}

		req := elastic.NewBulkIndexRequest().Index(c.indexName).Type(mappingType).Id(ID).Doc(products[i])
		bulkReq = bulkReq.Add(req)
	}

	_, err := bulkReq.Do(ctx)
	if err == nil {
		_, err = c.client.Flush().Index(c.indexName).Do(ctx)
	}
	return err
}

// Search performs a query
func (c *Client) Search(text string, count int) (*elastic.SearchResult, error) {
	ctx := context.Background()
	return c.client.Search().
		Index(c.indexName).
		Query(elastic.NewTermQuery("name", text)).
		From(0).
		Size(count).
		Do(ctx)
}

// Autocomplete performs a query
func (c *Client) Autocomplete(prefix string, count int) (*elastic.SearchResult, error) {
	ctx := context.Background()
	return c.client.Search().
		Index(c.indexName).
		Query(elastic.NewPrefixQuery("name_autocomplete", prefix)).
		From(0).
		Size(count).
		Do(ctx)
}

// Suggest performs a query
func (c *Client) Suggest(prefix string, count int) (*elastic.SearchResult, error) {
	ctx := context.Background()
	return c.client.Search().
		Index(c.indexName).
		Suggester(
			elastic.NewCompletionSuggester(c.indexName).
				Text(prefix).
				Field("suggestion"),
		).
		From(0).
		Size(count).
		Do(ctx)
}

// Delete a product
func (c *Client) Delete(sku int) error {
	ctx := context.Background()
	_, err := c.client.Delete().
		Index(c.indexName).
		Id(strconv.Itoa(sku)).
		Do(ctx)
	return err
}

// DeleteIndex remove an index
func (c *Client) DeleteIndex() error {
	ctx := context.Background()
	_, err := c.client.DeleteIndex(c.indexName).Do(ctx)
	return err
}

// NewClient creates an es client
func NewClient(hosts []string, login, password, indexName string) (*Client, error) {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(hosts...), elastic.SetBasicAuth(login, password), elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return nil, err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}
	return &Client{client, indexName}, nil
}
