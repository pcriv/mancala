package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/core"
	"github.com/pablocrivella/mancala/internal/service"
	"github.com/pablocrivella/mancala/internal/web/openapi"
)

// GamesHandler handles the requests for the games resource.
type GamesHandler struct {
	GameService service.GameService
}

func (h GamesHandler) Create(c echo.Context) error {
	b := new(openapi.CreateGameJSONRequestBody)
	err := c.Bind(b)
	if err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		return echo.NewHTTPError(code)
	}

	if strings.TrimSpace(b.Player1) == "" || strings.TrimSpace(b.Player2) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "players cannot be blank")
	}

	g, err := h.GameService.CreateGame(c.Request().Context(), b.Player1, b.Player2)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, g)
}

func (h GamesHandler) Show(c echo.Context, id string) error {
	g, err := h.GameService.FindGame(c.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrGameNotFound):
			return echo.NewHTTPError(http.StatusNotFound)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, g)
}
