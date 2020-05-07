package resources

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/engine"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
)

type (
	// GamesResource handles the requests for the games resource.
	GamesResource struct {
		GamesService *games.Service
	}

	CreateGameRequestBody struct {
		Player1 string `json:"player1"`
		Player2 string `json:"player2"`
	}

	UpdateGameRequestBody struct {
		PitIndex int `json:"pit_index"`
	}
)

func (h *GamesResource) Create(c echo.Context) error {
	b := new(CreateGameRequestBody)
	if err := c.Bind(b); err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		return echo.NewHTTPError(code)
	}

	if err := b.Validate(); err.Any() {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	g, err := h.GamesService.CreateGame(b.Player1, b.Player2)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, g)
}

func (h *GamesResource) Show(c echo.Context) error {
	g, err := h.GamesService.FindGame(c.Param("id"))
	if err != nil {
		switch e := err.(type) {
		case *persistence.ErrNotFound:
			return c.JSON(http.StatusNotFound, e.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, g)
}

func (h *GamesResource) Update(c echo.Context) error {
	b := new(UpdateGameRequestBody)
	if err := c.Bind(b); err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		return echo.NewHTTPError(code)
	}

	g, err := h.GamesService.ExecutePlay(c.Param("id"), b.PitIndex)
	if err != nil {
		switch e := err.(type) {
		case *persistence.ErrNotFound:
			return c.JSON(http.StatusNotFound, e.Error())
		case *engine.ErrInvalidPlay:
			return echo.NewHTTPError(
				http.StatusUnprocessableEntity,
				ValidationErrors{Errors: []ValidationError{{Field: "base", Msg: e.Error()}}},
			)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, e.Error())
		}
	}
	return c.JSON(http.StatusOK, g)
}

// Validate checks the request body for errors.
func (b CreateGameRequestBody) Validate() ValidationErrors {
	v := ValidationErrors{}
	if strings.TrimSpace(b.Player1) == "" {
		v.Add("player1", "cannot be blank")
	}

	if strings.TrimSpace(b.Player2) == "" {
		v.Add("player2", "cannot be blank")
	}
	return v
}
