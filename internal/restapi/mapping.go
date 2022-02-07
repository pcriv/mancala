package restapi

import (
	"github.com/pablocrivella/mancala/internal/engine"
	"github.com/pablocrivella/mancala/internal/openapi"
)

func OpenAPIGame(g engine.Game) openapi.Game {
	return openapi.Game{
		Id:         g.ID.String(),
		BoardSide1: OpenAPIBoard(g.BoardSide1),
		BoardSide2: OpenAPIBoard(g.BoardSide2),
		Result:     openapi.Result(g.Result),
		Turn:       openapi.Turn(g.Turn),
	}
}

func OpenAPIBoard(b engine.BoardSide) openapi.BoardSide {
	return openapi.BoardSide{
		Pits:   b.Pits[:],
		Store:  b.Store,
		Player: OpenAPIPlayer(b.Player),
	}
}

func OpenAPIPlayer(p engine.Player) openapi.Player {
	return openapi.Player{
		Name:  p.Name,
		Score: p.Score,
	}
}
