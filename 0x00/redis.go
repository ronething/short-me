package main

import (
	"github.com/go-redis/redis"
	"time"
)

const (
	URLIDKEY           = "next.url.id"
	ShortLinkKey       = "shortLink:%s:url"
	URLHashKey         = "urlHash:%s:url"
	ShortLinkDetailKey = "shortLink:%s:detail"
)

// redis cli
type RedisCli struct {
	Cli *redis.Client
}

type URLDetail struct {
	URL                 string        `json:"url"`
	CreateAt            string        `json:"create_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

// init redis client

func NewRedisCli(addr string, passwd string, db int) *RedisCli {
	var (
		c   *redis.Client
		err error
	)
	c = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	if _, err = c.Ping().Result(); err != nil {
		panic(err)
	}
	return &RedisCli{Cli: c}

}
