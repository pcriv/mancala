package handlers

import (
	"errors"
	"net/http"

	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/engine"
	"github.com/pablocrivella/mancala/persistence"
)

type (
	// Games handles the request for the games resource.
	Games struct {
		Repo persistence.Repo
	}

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

func (h *Games) Create(c echo.Context) error {
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

	if err := h.Repo.SaveGame(g); err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, g)
}

func (h *Games) Get(c echo.Context) error {
	id := c.Param("id")
	game, err := h.Repo.GetGame(id)

	if errors.Is(err, persistence.ErrNotFound) {
		return c.JSON(http.StatusNotFound, nil)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, game)
}

func (h *Games) Update(c echo.Context) error {
	p := new(playParams)

	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id := c.Param("id")
	game, err := h.Repo.GetGame(id)

	if errors.Is(err, persistence.ErrNotFound) {
		return c.JSON(http.StatusNotFound, nil)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = game.PlayTurn(p.PitIndex)

	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, validationErrors{[]string{err.Error()}})
	}

	h.Repo.SaveGame(*game)

	return c.JSON(http.StatusOK, game)
}
