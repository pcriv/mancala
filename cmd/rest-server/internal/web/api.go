package web

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/internal/openapi"
)

var _ openapi.ServerInterface = API{}

func NewAPI(spec *openapi3.T, svc mancala.Service) API {
	return API{
		DocsHandler:  DocsHandler{Spec: spec},
		GamesHandler: GamesHandler{GameService: svc},
		PlaysHandler: PlaysHandler{GameService: svc},
	}
}

type API struct {
	DocsHandler
	GamesHandler
	PlaysHandler
}
