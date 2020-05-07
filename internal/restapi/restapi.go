package restapi

import (
	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/restapi/resources"
)

type (
	Resource interface {
		Show(echo.Context) error
		Create(echo.Context) error
		Update(echo.Context) error
	}

	RequestBody interface {
		Validate() resources.ValidationErrors
	}

	App struct {
		GamesResource Resource
	}
)
