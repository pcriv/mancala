package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablocrivella/mancala/handlers"
	"github.com/pablocrivella/mancala/persistence"
)

func main() {
	e := echo.New()

	redisURL, ok := os.LookupEnv("REDIS_URL")

	if !ok {
		panic("missing env variable: REDIS_URL")
	}

	repo, err := persistence.CreateRepo(redisURL)

	if err != nil {
		panic(err)
	}

	e.File("/docs", "public/docs.html")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	h := &handlers.Games{Repo: repo}

	v1 := e.Group("/v1")
	v1.POST("/games", h.Create)
	v1.GET("/games/:id", h.Get)
	v1.PATCH("/games/:id", h.Update)

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
