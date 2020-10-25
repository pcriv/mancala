package resources

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/engine"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
	"github.com/pablocrivella/mancala/internal/pkg/openapi"
)

// GamesResource handles the requests for the games resource.
type GamesResource struct {
	GamesService games.Service
}

func (h GamesResource) Create(c echo.Context) error {
	b := new(openapi.CreateGameJSONBody)
	if err := c.Bind(b); err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		return NewErrorResponse(code)
	}

	errors := []string{}
	if strings.TrimSpace(b.Player1) == "" {
		errors = append(errors, "player1 cannot be blank")
	}
	if strings.TrimSpace(b.Player2) == "" {
		errors = append(errors, "player2 cannot be blank")
	}
	if len(errors) > 0 {
		return NewErrorResponse(http.StatusUnprocessableEntity, errors...)
	}

	g, err := h.GamesService.CreateGame(b.Player1, b.Player2)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, g)
}

func (h GamesResource) Show(c echo.Context, id string) error {
	g, err := h.GamesService.FindGame(id)
	if err != nil {
		switch e := err.(type) {
		case *persistence.NotFoundError:
			return NewErrorResponse(http.StatusNotFound, e.Error())
		default:
			return NewErrorResponse(http.StatusInternalServerError, e.Error())
		}
	}

	return c.JSON(http.StatusOK, g)
}

func (h GamesResource) Update(c echo.Context, id string) error {
	b := new(openapi.UpdateGameJSONBody)
	if err := c.Bind(b); err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		return NewErrorResponse(code)
	}

	g, err := h.GamesService.ExecutePlay(id, b.PitIndex)
	if err != nil {
		switch e := err.(type) {
		case *persistence.NotFoundError:
			return NewErrorResponse(http.StatusNotFound, e.Error())
		case *engine.InvalidPlayError:
			return NewErrorResponse(http.StatusUnprocessableEntity, e.Error())
		default:
			return NewErrorResponse(http.StatusInternalServerError, e.Error())
		}
	}
	return c.JSON(http.StatusOK, g)
}
