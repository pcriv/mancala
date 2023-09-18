package restapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
	"github.com/pablocrivella/mancala/internal/openapi"
)

type (
	// GamesHandler handles the requests for the games resource.
	GamesHandler struct {
		GamesService games.Service
	}
)

func (h GamesHandler) Create(c echo.Context) error {
	b := new(openapi.CreateGameJSONRequestBody)
	if err := c.Bind(b); err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		return echo.NewHTTPError(code)
	}

	if strings.TrimSpace(b.Player1) == "" || strings.TrimSpace(b.Player2) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "players cannot be blank")
	}

	g, err := h.GamesService.CreateGame(b.Player1, b.Player2)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, g)
}

func (h GamesHandler) Show(c echo.Context, id string) error {
	g, err := h.GamesService.FindGame(id)
	if err != nil {
		switch {
		case errors.Is(err, persistence.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, g)
}
