package restapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/pkg/openapi"
	"github.com/pablocrivella/mancala/internal/restapi/resources"
)

type Server struct {
	Games resources.GamesResource
}

func NewServer(gamesResources resources.GamesResource) Server {
	return Server{
		Games: gamesResources,
	}
}

func (s Server) CreateGame(ctx echo.Context) error {
	return s.Games.Create(ctx)
}

func (s Server) ShowGame(ctx echo.Context, id string) error {
	return s.Games.Show(ctx, id)
}

func (s Server) UpdateGame(ctx echo.Context, id string) error {
	return s.Games.Update(ctx, id)
}

func (s Server) ShowSwaggerSpec(ctx echo.Context) error {
	spec, err := openapi.GetSwagger()
	if err != nil {
		return resources.NewErrorResponse(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, spec)
}
