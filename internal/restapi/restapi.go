package restapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/openapi"
)

func NewServer(gamesService games.Service) Server {
	return Server{
		Games: GamesHandler{GamesService: gamesService},
		Plays: PlaysHandler{GamesService: gamesService},
	}
}

type Server struct {
	Games GamesHandler
	Plays PlaysHandler
}

func (s Server) CreateGame(ctx echo.Context) error {
	return s.Games.Create(ctx)
}

func (s Server) ShowGame(ctx echo.Context, id openapi.GameID) error {
	return s.Games.Show(ctx, string(id))
}

func (s Server) CreatePlay(ctx echo.Context, id openapi.GameID) error {
	return s.Plays.Create(ctx, string(id))
}

func (s Server) ShowSwaggerSpec(ctx echo.Context) error {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, spec)
}
