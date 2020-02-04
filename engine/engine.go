package engine

import (
	"github.com/google/uuid"
)

// NewGame initializes a manala game for the given players
func NewGame(player1Name string, player2Name string) Game {
	game := Game{}
	game.ID = uuid.New()
	game.BoardSide1.setup(player1Name)
	game.BoardSide2.setup(player2Name)
	game.Result = Undefined

	return game
}
