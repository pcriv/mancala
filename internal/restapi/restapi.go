package restapi

import (
	"github.com/pablocrivella/mancala/internal/restapi/resources"
)

type (
	RequestBody interface {
		Validate() resources.ValidationErrors
	}

	Resources struct {
		Games resources.GamesResource
	}
)
