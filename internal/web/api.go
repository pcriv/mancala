package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/service"
	"github.com/pablocrivella/mancala/internal/web/handlers"
	"github.com/pablocrivella/mancala/internal/web/openapi"
)

func NewAPI(gameService service.GameService) API {
	return API{
		Games: handlers.GamesHandler{GameService: gameService},
		Plays: handlers.PlaysHandler{GameService: gameService},
	}
}

type API struct {
	Games handlers.GamesHandler
	Plays handlers.PlaysHandler
}

func (s API) CreateGame(ctx echo.Context) error {
	return s.Games.Create(ctx)
}

func (s API) ShowGame(ctx echo.Context, id openapi.GameID) error {
	return s.Games.Show(ctx, id)
}

func (s API) CreatePlay(ctx echo.Context, id openapi.GameID) error {
	return s.Plays.Create(ctx, id)
}

func (s API) ShowSwaggerSpec(ctx echo.Context) error {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, spec)
}
