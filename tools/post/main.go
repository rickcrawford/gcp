package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/storage"

	"github.com/rickcrawford/gcp/common/models"
)

func DigestString(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func doPost(url string, products []models.Product) error {

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(products)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, &buf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
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
const bucketName = "typeahead-catalogs"
const url = "https://autocomplete-dot-typeahead-183622.appspot.com/products/batch"

func main() {
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

		err = doPost(url, products[start:end])
		if err != nil {
			log.Fatal(err)
		}

		log.Println(start, end, len(products))
	}
	log.Println("done!")

}
