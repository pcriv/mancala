package main

import (
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/xo/dburl"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/postgres"
	"github.com/pablocrivella/mancala/internal/restapi"
	"github.com/pablocrivella/mancala/internal/restapi/resources"
)

func main() {
	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("missing env variable: DATABASE_URL")
	}
	db, err := dburl.Open(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	gameRepo := postgres.NewGameRepo(db)

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
