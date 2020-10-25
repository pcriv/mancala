package resources

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/pkg/openapi"
)

func NewErrorResponse(code int, e ...string) *echo.HTTPError {
	if len(e) == 0 {
		e = []string{http.StatusText(code)}
	}
	return echo.NewHTTPError(
		code,
		openapi.ErrorResponse{
			Errors: e,
		},
	)
}
