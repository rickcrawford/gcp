package commands

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/rickcrawford/autocomplete.es/handlers"
	"github.com/rickcrawford/autocomplete.es/models"
	"github.com/rickcrawford/autocomplete.es/pubsub"
)

func start(sig <-chan os.Signal) bool {
	var exit bool
	log.Println("starting application")

	go func() {
		ctx := context.Background()

		redisURL, err := url.Parse(viper.GetString("redis-url"))
		if err != nil {
			log.Fatal("error connecting to redis", err)
		}

		redisPool := &redigo.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,

			Dial: func() (redigo.Conn, error) {
				c, err := redigo.Dial("tcp", redisURL.Host)
				if err != nil {
					return nil, err
				}
				return c, err
			},

			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}

		defer redisPool.Close()

		esClient, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(strings.Split(viper.GetString("elastic-url"), ",")...), elastic.SetBasicAuth("elastic", "changeme"))
		if err != nil {
			log.Fatal("error connecting to elastic", err, viper.GetString("elastic-url"))
		}

		// setup PubSub
		projectID := viper.GetString("project-id")
		topicName := viper.GetString("topic-name")
		subscriptionName := viper.GetString("subscription-name")

		pubSubClient, err := pubsub.NewClient(projectID, topicName, subscriptionName)
		if err != nil {
			log.Fatal("error starting pubsub client", err)
		}
		defer pubSubClient.Close()

		// setup Elastic Search
		go indexer(ctx, pubSubClient, esClient)

		clientArgs := models.ClientArgs{
			ES:        esClient,
			RedisPool: redisPool,
		}

		router := chi.NewRouter()
		router.Use(middleware.RealIP)
		router.Use(middleware.Recoverer)
		router.Use(middleware.DefaultCompress)

		router.Mount("/", handlers.GetRoutes(clientArgs))

		// // set the application namespaace, and appengine context
		// http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 	// Chi creates a copy of the request, so you need to register the context immediately or
		// 	// you will get an out of flight request context panic
		// 	ctx := appengine.NewContext(r)

		// 	// set the namespace. this should be specific to the logged in user context...
		// 	// for example we would probably want to set this based on the tennat
		// 	// ctx, _ = appengine.Namespace(ctx, "namespace")

		// 	router.ServeHTTP(w, r.WithContext(ctx))
		// }))

		port := viper.GetString("http-port")

		srvHTTP := &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}
		log.Println("started proxy http", port)
		if err := srvHTTP.ListenAndServe(); err != http.ErrServerClosed {
			log.Println("unexpected error from proxy", err)
			// Send a TERM signal
			killProcess()
		}
	}()

	switch <-sig {
	case syscall.SIGINT, syscall.SIGTERM:
		exit = true

	// case syscall.SIGHUP:
	default:
		log.Println("reload")

	}

	return exit
}

func killProcess() {
	// Send a TERM signal
	if p, err := os.FindProcess(os.Getpid()); err == nil {
		p.Signal(os.Kill)
	}
}
