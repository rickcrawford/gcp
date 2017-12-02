package main

import (
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"

	"github.com/rickcrawford/gcp/common/models"
)

func DigestString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func getData(bucketName, filename string) ([]byte, error) {
	tempName := DigestString(filename)
	tempFilename := "tmp/" + tempName
	if _, err := os.Stat(tempFilename); os.IsExist(err) {
		log.Println("exists", tempFilename)
		return ioutil.ReadFile(tempFilename)
	}

	ctx := context.Background()

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer storageClient.Close()

	// Creates a Bucket instance.
	bucket := storageClient.Bucket(bucketName)
	rdr, err := bucket.Object(filename).ReadCompressed(true).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	gzr, err := gzip.NewReader(rdr)
	if err != nil {
		return nil, err
	}
	defer gzr.Close()

	data, err := ioutil.ReadAll(gzr)
	if err != nil {
		return nil, err
	}

	os.Mkdir("tmp", 0777)
	ioutil.WriteFile(tempFilename, data, 0777)

	return data, nil

}

const projectID = "typeahead-183622"
const topicName = "updates"
const subscriptionName = "publish-subscription"
const bucketName = "typeahead-catalogs"

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create a new topic with the given name.
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
		topic = client.Topic(topicName)
	}
	defer topic.Stop()

	data, err := getData(bucketName, "bestbuy/products.json.gz")
	if err != nil {
		log.Fatalf("Failed to get products: %v", err)
	}

	var products []models.Product
	err = json.Unmarshal(data, &products)
	if err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}

	data, err = getData(bucketName, "bestbuy/content.json.gz")
	if err != nil {
		log.Fatalf("Failed to get products: %v", err)
	}

	var images []models.Image
	err = json.Unmarshal(data, &images)
	if err != nil {
		log.Fatalf("Failed to decode: %v", err)
	}

	log.Println("load images")
	imageSku := make(map[int]models.Image, len(images))
	for _, image := range images {
		imageSku[image.SKU] = image
	}
	log.Println("done!")

	log.Println("load products")
	for i := range products {
		if image, isPresent := imageSku[products[i].SKU]; isPresent {
			products[i].Content = image.Content
			products[i].Image = image.Image
		}
	}
	log.Println("done!")

	const max = 200
	log.Println("generate messages...")
	for start := 0; start <= len(products); start = start + max {
		end := start + max
		if end > len(products) {
			end = len(products)
		}

		data, err = json.Marshal(models.Message{Products: products[start:end], Type: models.MessageTypeUpdate})
		if err != nil {
			log.Fatal(err)
		}

		res := topic.Publish(ctx, &pubsub.Message{Data: data})
		id, err := res.Get(ctx)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(start, end, len(products), id)
	}
	log.Println("done!")

	// id, err := res.Get(ctx)
	// if err != nil {
	// 	log.Fatal("error publishing message", err)
	// }
	// log.Printf("Published a message with a message ID: %s\n", id)

	// <-time.After(time.Second * 10)

	// err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
	// 	log.Printf("Got message: %s", m.Data)
	// 	m.Ack()
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
