package web

import (
	"errors"
	"net/http"


	"github.com/pcriv/mancala/cmd/rest-server/internal/openapimap"

	"github.com/labstack/echo/v4"
	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/internal/openapi"
)

// GamesHandler handles the requests for the games resource.
type GamesHandler struct {
	GameService mancala.Service
}

func (h GamesHandler) CreateGame(c echo.Context) error {
	b := new(openapi.CreateGameJSONRequestBody)
	err := c.Bind(b)
	if err != nil {
		code := http.StatusInternalServerError
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
		}
		return echo.NewHTTPError(code)
	}

	g, err := h.GameService.CreateGame(c.Request().Context(), b.Player1, b.Player2)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, openapimap.Game(g))
}

func (h GamesHandler) ShowGame(c echo.Context, id openapi.GameID) error {
	g, err := h.GameService.FindGame(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, mancala.ErrGameNotFound):
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, openapimap.Game(g))
}
