package persistence

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// NewRedisClient returns a redis client.
func NewRedisClient(url string) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
