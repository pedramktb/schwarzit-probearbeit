package redis

import (
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

func create(redisURI string) *redis.Client {
	addr, username, password := parseURI(redisURI)
	cache := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
	})

	if _, err := cache.Ping(cache.Context()).Result(); err != nil {
		panic(errors.Wrap(err, "failed to connect to cache"))
	}

	return cache
}

// parsing redis://username:password@addr
func parseURI(uri string) (addr, username, password string) {
	uri, _ = strings.CutPrefix(uri, "redis://")
	lastAds := strings.LastIndex(uri, "@")
	if lastAds == -1 {
		return addr, "", ""
	}
	firstColon := strings.Index(uri, ":")
	if firstColon == -1 {
		return addr, "", ""
	}
	return uri[lastAds+1:], uri[:firstColon], uri[firstColon+1 : lastAds]
}
