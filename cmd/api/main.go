package main

import (
	"errors"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
	"github.com/pablocrivella/mancala/internal/pkg/openapi"
	"github.com/pablocrivella/mancala/internal/restapi"
	"github.com/pablocrivella/mancala/internal/restapi/resources"
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
	redisClient, err := persistence.NewRedisClient(redisURL)
	if err != nil {
		return err
	}
	gameRepo := persistence.NewGameRepo(redisClient)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(restapi.SwaggerUI())

	s := restapi.NewServer(
		resources.GamesResource{
			GamesService: games.NewService(gameRepo),
		},
	)
	e.GET("swagger.json", s.ShowSwaggerSpec)
	openapi.RegisterHandlers(e.Group("/v1"), &s)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(":" + port))
	return nil
}
