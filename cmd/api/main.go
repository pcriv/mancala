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

	g, err := persistence.NewGameRepo(redisURL)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.File("/", "website/public/index.html")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	api := restapi.App{
		GamesResource: &resources.GamesResource{
			GamesService: games.NewService(g),
		},
	}
	v1 := e.Group("/v1")
	v1.GET("/games/:id", api.GamesResource.Show)
	v1.POST("/games", api.GamesResource.Create)
	v1.PATCH("/games/:id", api.GamesResource.Update)

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
