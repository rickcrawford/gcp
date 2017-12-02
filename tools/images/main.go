package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nfnt/resize"
	"github.com/rickcrawford/gcp/common/models"
)

var urlContent map[string]string
var lock sync.Mutex

func init() {
	urlContent = make(map[string]string)
}

func getImage(SKU int, URL string, depth int) (*models.Image, error) {
	log.Println(SKU, URL, depth)
	if depth == 2 {
		return nil, errors.New("too deep")
	}
	lock.Lock()
	content := urlContent[URL]
	lock.Unlock()

	if content != "" {
		return &models.Image{
			SKU:     SKU,
			Content: content,
			Image:   URL,
		}, nil
	}

	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return getImage(SKU, "http://img.bbystatic.com/BestBuy_US/images/products/nonsku/default_hardlines_m.gif", depth+1)
	}

	out := &models.Image{
		SKU:   SKU,
		Image: URL,
	}

	var contentType string
	var img image.Image
	if strings.HasSuffix(URL, ".gif") {
		img, err = gif.Decode(response.Body)
		contentType = "image/gif"
	} else {
		img, err = gif.Decode(response.Body)
		contentType = "image/jpeg"

	}
	if err != nil {
		return nil, err
	}
	img = resize.Thumbnail(50, 50, img, resize.Lanczos3)
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, nil)
	str := base64.StdEncoding.EncodeToString(buf.Bytes())
	out.Content = "data:" + contentType + ";base64," + str
	lock.Lock()
	urlContent[URL] = out.Content
	lock.Unlock()

	<-time.After(time.Second * 2)
	return out, nil
}

func main() {

	files, err := ioutil.ReadDir("./temp")
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	images := make([]models.Image, 0)

	for _, f := range files {
		rdr, _ := os.Open("./temp/" + f.Name())
		gzr, _ := gzip.NewReader(rdr)
		decoder := json.NewDecoder(gzr)

		var update models.Message
		decoder.Decode(&update)
		gzr.Close()
		rdr.Close()
		for _, product := range update.Products {
			if product.Content == "" {
				// image, err := getImage(product.SKU, product.Image, 0)
				// if err != nil {
				// 	log.Println(product.SKU, err)
				// } else {
				// 	encoder.Encode(image)
				// }
				product.Content = "data:image/gif;base64,/9j/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIADIAMgMBIgACEQEDEQH/xAGiAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgsQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+gEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoLEQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2gAMAwEAAhEDEQA/APf6KKCQASegoAKK4HUvGN+b1hZgQQrwFdAS3uc9Kji8bamn347eT6qR/I0AehUVl6Bqkur6cbmWJY28wphSSDjFalABRRRQAUjKGUqehGKWq19cNbWwkR41O9FzJ05YD1HrQBwOp+F9TgvZPJikuojysgOSR6H3FZkmk6jF9+xuV/7ZmvSX1Uxly1uRGpclt4+6jbWOPr2qEa9GVJERyr7Wy4A7dCep5HFADPCkDQeH4FdSrMzMQwwepraqnY3633m7Y3TYcfN3/wAD7VcoAKKKKAMttTaPe7mLIOPJJ2snOAWPYd84/Oj+1WkGRaZUnA3OBkhdx4x6dPWtTA9OtFAGYmq7hvEDNGBksDyMuVGB36etMFxFe2k0hgi3W7blBfK5wDk4+vQ1rVHDBFbqViQKCcn3NAGbHqoSKFmjVmlJyy4XPOOBk7j9CfrT11YyIDHBuYvtA8wYA2luTjgjbyO2RWngccdKMUAMikWaFJVztdQwz1wafRRQAUUUUAFFFFABRRRQAUUUUAf/2Q=="
				product.Image = "http://img.bbystatic.com/BestBuy_US/images/products/nonsku/default_hardlines_m.gif"
			}
			images = append(images, models.Image{
				SKU:     product.SKU,
				Image:   product.Image,
				Content: product.Content,
			})
		}
	}

	encoder.Encode(images)
	ioutil.WriteFile("content.json", buf.Bytes(), 0644)
}
