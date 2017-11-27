package models

import (
	redigo "github.com/garyburd/redigo/redis"
	elastic "gopkg.in/olivere/elastic.v5"
)

// ClientArgs are used to run this application
type ClientArgs struct {
	RedisPool *redigo.Pool
	ES        *elastic.Client
}
