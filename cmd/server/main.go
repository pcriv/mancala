package main

import (
	"context"
	"errors"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablocrivella/mancala/internal/service"
	"github.com/pablocrivella/mancala/internal/web"
	"github.com/pablocrivella/mancala/internal/web/openapi"
	"github.com/redis/go-redis/v9"

	redisstore "github.com/pablocrivella/mancala/internal/store/redis"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		return errors.New("missing env variable: REDIS_URL")
	}
	redisClient, err := newRedisClient(redisURL)
	if err != nil {
		return err
	}
	gameStore := redisstore.NewGameStore(redisClient)

	e := echo.New()
	e.File("/docs", "website/public/index.html")

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())

	s := web.NewAPI(service.NewGameService(gameStore))
	e.GET("swagger.json", s.ShowSwaggerSpec)
	openapi.RegisterHandlers(e.Group("/v1"), &s)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(":" + port))

	return nil
}

func newRedisClient(url string) (*redis.Client, error) {
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
