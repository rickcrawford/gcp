package commands

import (
	"fmt"
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

	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/elastic"
	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/handlers"
	"github.com/rickcrawford/gcp/kubernetes/autocomplete.es/pubsub"
)

var envRedisTarget string
var envElasticTarget string

func init() {
	redisHost := os.Getenv("REDIS_SERVICE_HOST")
	redisPort := os.Getenv("REDIS_SERVICE_PORT")
	envRedisTarget = fmt.Sprintf("redis://%s:%s", redisHost, redisPort)

	elasticHost := os.Getenv("ELASTICSEARCH_SERVICE_HOST")
	elasticPort := os.Getenv("ELASTICSEARCH_SERVICE_PORT")
	envElasticTarget = fmt.Sprintf("http://%s:%s", elasticHost, elasticPort)
}

func start(sig <-chan os.Signal) bool {
	var exit bool
	log.Println("starting application")

	go func() {

		redisTarget := viper.GetString("redis-url")
		if redisTarget == "" {
			redisTarget = envRedisTarget
		}

		redisURL, err := url.Parse(redisTarget)
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

		elasticTarget := viper.GetString("elastic-url")
		if elasticTarget == "" {
			elasticTarget = envElasticTarget
		}

		elasticHosts := strings.Split(elasticTarget, ",")
		elasticLogin := viper.GetString("elastic-login")
		elasticPassword := viper.GetString("elastic-password")
		indexName := viper.GetString("elastic-index-name")
		debug := viper.GetBool("debug")

		esClient, err := elastic.NewClient(elasticHosts, elasticLogin, elasticPassword, indexName, debug)
		if err != nil {
			log.Fatal("error loading search", err)
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

		router := chi.NewRouter()
		router.Use(middleware.RealIP)
		router.Use(middleware.Recoverer)
		router.Use(middleware.DefaultCompress)

		router.Mount("/", handlers.GetRoutes(esClient, pubSubClient, redisPool))

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
