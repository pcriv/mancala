package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
	"github.com/pablocrivella/mancala/internal/restapi"
	"github.com/pablocrivella/mancala/internal/restapi/resources"
)

func main() {
	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		panic("missing env variable: REDIS_URL")
	}
	redisClient, err := persistence.NewRedisClient(redisURL)
	if err != nil {
		panic(err)
	}
	gameRepo := persistence.NewGameRepo(redisClient)

	e := echo.New()
	e.File("/", "website/public/index.html")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	r := restapi.Resources{
		Games: resources.GamesResource{
			GamesService: games.NewService(gameRepo),
		},
	}
	v1 := e.Group("/v1")
	v1.GET("/games/:id", r.Games.Show)
	v1.POST("/games", r.Games.Create)
	v1.PATCH("/games/:id", r.Games.Update)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
