package persistence

import "github.com/go-redis/redis/v7"

// NewRedisClient returns a redis client.
func NewRedisClient(url string) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
