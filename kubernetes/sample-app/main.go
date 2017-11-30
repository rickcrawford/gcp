package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/olivere/elastic"
)

func main() {

	// redisHost := os.Getenv("REDIS_SERVICE_HOST")
	// redisPort := os.Getenv("REDIS_SERVICE_PORT")
	// redisTarget := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// elasticHost := os.Getenv("ELASTICSEARCH_SERVICE_HOST")
	// elasticPort := os.Getenv("ELASTICSEARCH_SERVICE_PORT")
	// elasticTarget := fmt.Sprintf("%s:%s", elasticHost, elasticPort)

	http.HandleFunc("/status", func(rw http.ResponseWriter, req *http.Request) {
		for _, key := range os.Environ() {
			fmt.Fprintln(rw, key)
		}

		conn, err := redigo.Dial("tcp", "redis:6379")
		if err != nil {
			fmt.Fprintln(rw, "redis error", err)
		} else {

			_, err = conn.Do("PING")
			if err != nil {
				fmt.Fprintln(rw, "ping error", err)
			}
			conn.Close()
		}
		// test

		hosts := []string{"http://" + os.Getenv("ELASTICSEARCH_PORT_9200_TCP_ADDR") + ":9200"}
		login := "elastic"
		password := "changeme"
		options := make([]elastic.ClientOptionFunc, 0)
		options = append(options, elastic.SetSniff(false), elastic.SetURL(hosts...))
		if login != "" {
			options = append(options, elastic.SetBasicAuth(login, password))
		}
		options = append(options, elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))

		client, err := elastic.NewSimpleClient(elastic.SetURL("http://elasticsearch:9200"))
		if err != nil {
			fmt.Fprintln(rw, "elastic error", err)
		} else {

			info, err := client.ClusterHealth().Do(context.Background())
			if err != nil {
				fmt.Fprintln(rw, "elastic info error", err)
			}
			fmt.Fprintln(rw, "cluster name", info)
		}

	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
