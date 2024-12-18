package web

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

type DocsHandler struct {
	Spec *openapi3.T
}

func (s DocsHandler) RenderSwagger(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, s.Spec)
}

func (s DocsHandler) RenderDocs(ctx echo.Context) error {
	return ctx.File("website/public/index.html")
}
