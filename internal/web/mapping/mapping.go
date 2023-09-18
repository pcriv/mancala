package mapping

import (
	"github.com/pablocrivella/mancala/internal/core"
	"github.com/pablocrivella/mancala/internal/web/openapi"
)

func ToOpenAPIGame(g core.Game) openapi.Game {
	return openapi.Game{
		Id:         g.ID.String(),
		BoardSide1: ToOpenAPIBoard(g.BoardSide1),
		BoardSide2: ToOpenAPIBoard(g.BoardSide2),
		Result:     openapi.Result(g.Result),
		Turn:       openapi.Turn(g.Turn),
	}
}

func ToOpenAPIBoard(b core.BoardSide) openapi.BoardSide {
	return openapi.BoardSide{
		Pits:   b.Pits[:],
		Store:  b.Store,
		Player: ToOpenAPIPlayer(b.Player),
	}
}

func ToOpenAPIPlayer(p core.Player) openapi.Player {
	return openapi.Player{
		Name:  p.Name,
		Score: p.Score,
	}
}
