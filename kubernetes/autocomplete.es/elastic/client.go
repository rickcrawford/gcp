package elastic

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/rickcrawford/gcp/common/models"
	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/pubsub"
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
					"index":false
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
					"index":false
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
				"keywords":{
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
	client       *elastic.Client
	pubSubClient *pubsub.Client

	indexName string
	exit      chan interface{}
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
	bulkReq := c.client.Bulk().Index(c.indexName).Type(mappingType)

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
		Type(mappingType).
		Id(strconv.Itoa(sku)).
		Do(ctx)
	return err
}

// BulkDelete products
func (c *Client) BulkDelete(skus []int) error {
	ctx := context.Background()
	bulkReq := c.client.Bulk().Index(c.indexName).Type(mappingType)

	for i := range skus {
		req := elastic.NewBulkDeleteRequest().Index(c.indexName).Id(strconv.Itoa(skus[i]))
		bulkReq = bulkReq.Add(req)
	}

	_, err := bulkReq.Do(ctx)
	if err == nil {
		_, err = c.client.Flush().Index(c.indexName).Do(ctx)
	}
	return err
}

// DeleteIndex remove an index
func (c *Client) DeleteIndex() error {
	ctx := context.Background()
	_, err := c.client.DeleteIndex(c.indexName).Do(ctx)
	return err
}

// Close exits any processes
func (c *Client) Close() error {
	close(c.exit)
	return nil
}

// updateProcessor will read messages off of a queue and process them.
func (c *Client) updateProcessor() {
	log.Println("--- starting update processor ---")

	var update *models.Message
	var err error
	for {
		select {
		case data := <-c.pubSubClient.GetProductUpdate():
			log.Println("--- processing an update ---")

			update = new(models.Message)
			if err = json.Unmarshal(data, update); err == nil {

				switch update.Type {
				case models.MessageTypeDelete:
					skus := make([]int, len(update.Products))
					for i := range update.Products {
						skus[i] = update.Products[i].SKU
					}
					err = c.BulkDelete(skus)
				case models.MessageTypeUpdate:
					err = c.BulkIndex(update.Products)
				}
			}

			log.Println("update!")

			if err != nil {
				log.Println("error processing update", err)
			}

		case <-c.exit:
			return
		}

	}
}

// NewClient creates an es client
func NewClient(hosts []string, login, password, indexName string, debug bool, pubSubClient *pubsub.Client) (*Client, error) {
	log.Println("starting elastic", hosts, login, password, indexName)

	options := make([]elastic.ClientOptionFunc, 0)
	options = append(options, elastic.SetSniff(false), elastic.SetURL(hosts...))
	if login != "" {
		options = append(options, elastic.SetBasicAuth(login, password))
	}
	if debug {
		options = append(options, elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	}

	client, err := elastic.NewClient(options...)
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

	esClient := &Client{
		client:       client,
		indexName:    indexName,
		pubSubClient: pubSubClient,
		exit:         make(chan interface{}),
	}

	if pubSubClient != nil {
		go esClient.updateProcessor()
	}

	log.Println("done!")
	return esClient, nil
}
