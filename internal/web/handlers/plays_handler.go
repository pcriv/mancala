package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/core"
	"github.com/pablocrivella/mancala/internal/service"
	"github.com/pablocrivella/mancala/internal/web/mapping"
	"github.com/pablocrivella/mancala/internal/web/openapi"
)

type PlaysHandler struct {
	GameService service.GameService
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

	g, err := h.GameService.ExecutePlay(c.Request().Context(), gameID, b.PitIndex)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrGameNotFound):
			return echo.NewHTTPError(http.StatusNotFound)
		case errors.Is(err, core.ErrInvalidPlay):
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}
	return c.JSON(http.StatusOK, openapi.PlayCreated{
		PlayedPitIndex: b.PitIndex,
		Game:           mapping.ToOpenAPIGame(g),
	})
}
