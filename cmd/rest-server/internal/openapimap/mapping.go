package openapimap

import (
	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/internal/openapi"
)

func Game(g mancala.Game) openapi.Game {
	return openapi.Game{
		Id:         g.ID.String(),
		BoardSide1: Board(g.BoardSide1),
		BoardSide2: Board(g.BoardSide2),
		Result:     openapi.Result(g.Result),
		Turn:       openapi.Turn(g.Turn),
	}
}

func Board(b mancala.BoardSide) openapi.BoardSide {
	return openapi.BoardSide{
		Pits:   b.Pits[:],
		Store:  b.Store,
		Player: Player(b.Player),
	}
}

func Player(p mancala.Player) openapi.Player {
	return openapi.Player{
		Name:  p.Name,
		Score: p.Score,
	}
}
