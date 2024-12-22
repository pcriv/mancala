package mancala

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	TurnPlayer1 Turn = "player1"
	TurnPlayer2 Turn = "player2"
)

const (
	// Undefined represents a yet undefined result for a game
	Undefined Result = "Undefined"
	// Player1Wins represents a game result where Player1 wins the game
	Player1Wins Result = "Player1Wins"
	// Player2Wins represents a game result where Player2 wins the game
	Player2Wins Result = "Player2Wins"
	// Tie represents a game result where both players have the same score
	Tie Result = "Tie"
)

const (
	pitsPerSide  int64 = 6
	stonesPerPit int64 = 6
)

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
	Turn string

	// Result represents the result of the game after it is done, can be Player1Wins, Player2Wins or Tie
	Result string
)

// NewGame initializes a mancala game for the given players
func NewGame(player1Name string, player2Name string) Game {
	game := Game{}
	game.ID = uuid.New()
	game.BoardSide1.setup(player1Name)
	game.BoardSide2.setup(player2Name)
	game.Turn = TurnPlayer1
	game.Result = Undefined

	return game
}

// PlayTurn from the given pitIndex for the current playingSide.
func (game *Game) PlayTurn(pitIndex int64) error {
	if pitIndex < 0 || pitIndex > pitsPerSide-1 {
		return fmt.Errorf("%w: pit index is invalid", ErrInvalidPlay)
	}

	if game.IsDone() {
		return fmt.Errorf("%w: game is already done", ErrInvalidPlay)
	}

	stones := game.playingSide().pickStones(pitIndex)
	if stones == 0 {
		return fmt.Errorf("%w: selected pit is empty", ErrInvalidPlay)
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
	if game.Turn == TurnPlayer1 {
		game.Turn = TurnPlayer2
	} else {
		game.Turn = TurnPlayer1
	}
}

func (game *Game) playingSide() *BoardSide {
	if game.Turn == TurnPlayer1 {
		return &game.BoardSide1
	}

	return &game.BoardSide2
}

func (game *Game) opposingSide() *BoardSide {
	if game.Turn == TurnPlayer1 {
		return &game.BoardSide2
	}

	return &game.BoardSide1
}
