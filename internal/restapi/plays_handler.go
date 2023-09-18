package restapi

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/engine"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
	"github.com/pablocrivella/mancala/internal/openapi"
)

type PlaysHandler struct {
	GamesService games.Service
}

func (h PlaysHandler) Create(c echo.Context, gameID string) error {
	b := new(openapi.CreatePlayJSONRequestBody)
	if err := c.Bind(b); err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		return echo.NewHTTPError(code)
	}

	g, err := h.GamesService.ExecutePlay(gameID, b.PitIndex)
	if err != nil {
		switch {
		case errors.Is(err, persistence.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound)
		case errors.Is(err, engine.ErrInvalidPlay):
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return c.JSON(http.StatusOK, openapi.PlayCreated{
		PlayedPitIndex: b.PitIndex,
		Game:           OpenAPIGame(g),
	})
}
