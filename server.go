package main

import (
	"net/http"
	"os"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pablocrivella/mancala/engine"
	"github.com/pablocrivella/mancala/repo"
	"golang.org/x/xerrors"
)

type (
	newGameParams struct {
		Player1 string `json:"player1"`
		Player2 string `json:"player2"`
	}

	playParams struct {
		PitIndex int `json:"pit_index"`
	}

	validationErrors struct {
		Errors []string `json:"errors"`
	}
)

func main() {
	e := echo.New()

	e.File("/docs", "public/docs.html")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	v1 := e.Group("/v1")
	v1.POST("/games", createGame)
	v1.GET("/games/:id", getGame)
	v1.PATCH("/games/:id", updateGame)

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func createGame(c echo.Context) error {
	p := new(newGameParams)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	e := validate.Validate(
		&validators.StringIsPresent{Field: p.Player1, Name: "player1"},
		&validators.StringIsPresent{Field: p.Player2, Name: "player2"},
	)

	if len(e.Errors) != 0 {
		v := validationErrors{}

		for _, errors := range e.Errors {
			v.Errors = append(v.Errors, errors...)
		}

		return c.JSON(http.StatusUnprocessableEntity, v)
	}

	g := engine.NewGame(p.Player1, p.Player2)

	if err := repo.SaveGame(g); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, g)
}

func getGame(c echo.Context) error {
	id := c.Param("id")
	game, err := repo.GetGame(id)

	if xerrors.Is(err, repo.ErrNotFound) {
		return c.JSON(http.StatusNotFound, nil)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, game)
}

func updateGame(c echo.Context) error {
	p := new(playParams)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id := c.Param("id")
	game, err := repo.GetGame(id)

	if xerrors.Is(err, repo.ErrNotFound) {
		return c.JSON(http.StatusNotFound, nil)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = game.PlayTurn(p.PitIndex)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrors{[]string{err.Error()}})
	}

	repo.SaveGame(*game)

	return c.JSON(http.StatusOK, game)
}
