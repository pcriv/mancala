package web

import (
	"errors"
	"net/http"


	"github.com/pcriv/mancala/cmd/rest-server/internal/openapimap"

	"github.com/pcriv/mancala/internal/mancala"

	"github.com/labstack/echo/v4"
	"github.com/pcriv/mancala/internal/openapi"
)

type PlaysHandler struct {
	GameService mancala.Service
}

func (h PlaysHandler) CreatePlay(c echo.Context, gameID openapi.GameID) error {
	b := new(openapi.CreatePlayJSONRequestBody)
	if err := c.Bind(b); err != nil {
		code := http.StatusInternalServerError
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
		}
		return echo.NewHTTPError(code)
	}

	g, err := h.GameService.ExecutePlay(c.Request().Context(), gameID, b.PitIndex)
	if err != nil {
		switch {
		case errors.Is(err, mancala.ErrGameNotFound):
			return echo.NewHTTPError(http.StatusNotFound)
		case errors.Is(err, mancala.ErrInvalidPlay):
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, openapi.PlayCreated{
		PlayedPitIndex: b.PitIndex,
		Game:           openapimap.Game(g),
	})
}
