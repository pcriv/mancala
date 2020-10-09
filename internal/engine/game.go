package engine

import (
	"github.com/google/uuid"
)

const (
	// Player1Turn represents the turn of the Player #1
	Player1Turn Turn = iota
	// Player2Turn represents the turn of the Player #2
	Player2Turn
)

const (
	// Undefined represents a yet undefined result for a game
	Undefined Result = iota
	// Player1Wins represents a game result where Player1 wins the game
	Player1Wins
	// Player2Wins represents a game result where Player2 wins the game
	Player2Wins
	// Tie represents a game result where both players have the same score
	Tie
)

const pitsPerSide int64 = 6

type (
	// Game represents a Mancala game.
	Game struct {
		ID         uuid.UUID `json:"id"`
		Turn       Turn      `json:"turn"`
		BoardSide1 BoardSide `json:"board_side1"`
		BoardSide2 BoardSide `json:"board_side2"`
		Result     Result    `json:"result"`
	}

	// Turn represents a turn on the mancala game, can be Player1Turn or Player2Turn
	Turn int64

	// Result represents the result of the game after it is done, can be Player1Wins, Player2Wins or Tie
	Result int64
)

// PlayTurn from the given pitIndex for the current playingSide.
func (game *Game) PlayTurn(pitIndex int64) error {
	if pitIndex < 0 || pitIndex > pitsPerSide-1 {
		return &InvalidPlayError{Msg: "pit index is invalid"}
	}

	if game.IsDone() {
		return &InvalidPlayError{Msg: "game is already done"}
	}

	stones := game.playingSide().pickStones(pitIndex)

	if stones == 0 {
		return &InvalidPlayError{Msg: "selected pit is empty"}
	}

	game.placeStones(pitIndex+1, stones)

	if game.IsDone() {
		game.calculateScores()
	}

	return nil
}

// IsDone returns true when the Game is finished.
func (game *Game) IsDone() bool {
	return game.BoardSide1.isEmpty() || game.BoardSide2.isEmpty()
}

func (game *Game) calculateScores() {
	game.BoardSide1.score()
	game.BoardSide2.score()

	if game.BoardSide1.Player.Score > game.BoardSide2.Player.Score {
		game.Result = Player1Wins
	} else if game.BoardSide2.Player.Score > game.BoardSide1.Player.Score {
		game.Result = Player2Wins
	} else {
		game.Result = Tie
	}
}

func (game *Game) placeStones(pitIndex int64, stones int64) {
	for index := pitIndex; index < pitsPerSide; index++ {
		stones = game.playingSide().placeStone(index, stones)

		if stones == 0 {
			if game.playingSide().Pits[index] == 1 {
				stones = game.playingSide().pickStones(index)

				oppositePitIndex := (index - pitsPerSide + 1) * -1
				oppositePitStones := game.opposingSide().pickStones(oppositePitIndex)

				stones += oppositePitStones

				game.playingSide().captureAllStones(stones)
			}

			game.changeTurn()

			return
		}
	}

	if stones != 0 {
		stones = game.playingSide().captureStone(stones)
	}

	if stones == 0 {
		return
	}

	for index := 0; int64(index) < pitsPerSide; index++ {
		stones = game.opposingSide().placeStone(int64(index), stones)

		if stones == 0 {
			game.changeTurn()

			return
		}
	}

	game.placeStones(0, stones)
}

func (game *Game) changeTurn() {
	game.Turn = 1 - game.Turn
}

func (game *Game) playingSide() *BoardSide {
	if game.Turn == Player1Turn {
		return &game.BoardSide1
	}

	return &game.BoardSide2
}

func (game *Game) opposingSide() *BoardSide {
	if game.Turn == Player1Turn {
		return &game.BoardSide2
	}

	return &game.BoardSide1
}
