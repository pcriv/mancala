package mancala

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	game := NewGame("Rick", "Morty")

	assert.Equal(t, game.BoardSide1.Pits, fullSidePits)
	assert.Equal(t, game.BoardSide2.Pits, fullSidePits)
	assert.Equal(t, "Rick", game.BoardSide1.Player.Name)
	assert.Equal(t, "Morty", game.BoardSide2.Player.Name)
	assert.Equal(t, TurnPlayer1, game.Turn)
	assert.Equal(t, Undefined, game.Result)
}

func TestGame_PlayTurn(t *testing.T) {
	testCases := []struct {
		name     string
		game     Game
		pitIndex int64
		wantGame Game
		wantErr  string
	}{
		{
			name:     "when player1 plays from pit 0 and turn passes to player2",
			pitIndex: 0,
			game:     NewGame("Rick", "Morty"),
			wantGame: Game{
				BoardSide1: BoardSide{
					Pits:  buildPits(0, 5, 5, 5, 5, 4),
					Store: 0,
				},
				BoardSide2: BoardSide{
					Pits:  fullSidePits,
					Store: 0,
				},
				Result: Undefined,
				Turn:   TurnPlayer2,
			},
		},
		{
			name:     "when player1 moves the last stone to an empty own pit",
			pitIndex: 0,
			game: Game{
				BoardSide1: BoardSide{
					Pits: buildPits(1, 0, 0, 0, 0, 1),
				},
				BoardSide2: BoardSide{
					Pits: buildPits(6, 6, 6, 6, 1, 6),
				},
				Result: Undefined,
				Turn:   TurnPlayer1,
			},
			wantGame: Game{
				BoardSide1: BoardSide{
					Pits:  buildPits(0, 0, 0, 0, 0, 1),
					Store: 2,
				},
				BoardSide2: BoardSide{
					Pits:  buildPits(6, 6, 6, 6, 0, 6),
					Store: 0,
				},
				Result: Undefined,
				Turn:   TurnPlayer2,
			},
		},
		{
			name:     "when player1 places stones on player2 pit but skipping player2 store",
			pitIndex: 5,
			game: Game{
				BoardSide1: BoardSide{
					Pits: buildPits(1, 1, 0, 0, 0, 8),
				},
				BoardSide2: BoardSide{
					Pits: fullSidePits,
				},
				Result: Undefined,
				Turn:   TurnPlayer1,
			},
			wantGame: Game{
				BoardSide1: BoardSide{
					Pits:  buildPits(2, 1, 0, 0, 0, 0),
					Store: 1,
				},
				BoardSide2: BoardSide{
					Pits:  buildPits(5, 5, 5, 5, 5, 5),
					Store: 0,
				},
				Result: Undefined,
				Turn:   TurnPlayer2,
			},
		},
		{
			name:     "when it is the last turn before the end of the game",
			pitIndex: 5,
			game: Game{
				BoardSide1: BoardSide{
					Pits: buildPits(0, 0, 0, 0, 0, 1),
				},
				BoardSide2: BoardSide{
					Pits: buildPits(1, 1, 1, 1, 1, 1),
				},
				Result: Undefined,
				Turn:   TurnPlayer1,
			},
			wantGame: Game{
				BoardSide1: BoardSide{
					Pits:  buildPits(0, 0, 0, 0, 0, 0),
					Store: 1,
				},
				BoardSide2: BoardSide{
					Pits:  buildPits(1, 1, 1, 1, 1, 1),
					Store: 6,
				},
				Turn:   TurnPlayer1,
				Result: Player2Wins,
			},
		},
		{
			name:     "when the pit index is invalid",
			pitIndex: 8,
			game:     Game{},
			wantGame: Game{},
			wantErr:  "invalid play: pit index is invalid",
		},
		{
			name:     "when the pit index is invalid",
			pitIndex: 0,
			game: Game{
				BoardSide1: BoardSide{
					Pits: buildPits(0, 0, 0, 0, 0, 1),
				},
				BoardSide2: BoardSide{
					Pits: fullSidePits,
				},
				Result: Undefined,
				Turn:   TurnPlayer1,
			},
			wantGame: Game{},
			wantErr:  "invalid play: selected pit is empty",
		},
		{
			name:     "the game is already done",
			pitIndex: 0,
			game: Game{
				BoardSide1: BoardSide{
					Pits: buildPits(0, 0, 0, 0, 0, 1),
				},
			},
			wantGame: Game{},
			wantErr:  "invalid play: game is already done",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.game.PlayTurn(tc.pitIndex)
			if tc.wantErr == "" {
				assert.Equal(t, tc.wantGame.BoardSide1.Pits, tc.game.BoardSide1.Pits)
				assert.Equal(t, tc.wantGame.BoardSide2.Pits, tc.game.BoardSide2.Pits)
				assert.Equal(t, tc.wantGame.BoardSide1.Store, tc.game.BoardSide1.Store)
				assert.Equal(t, tc.wantGame.BoardSide2.Store, tc.game.BoardSide2.Store)
				assert.Equal(t, tc.wantGame.Turn, tc.game.Turn)
				assert.Equal(t, tc.wantGame.Result, tc.game.Result)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
		})
	}
}

func TestIsDone(t *testing.T) {
	testCases := []struct {
		name string
		game Game
		want bool
	}{
		{
			name: "when one of the sides is empty",
			game: Game{
				BoardSide1: BoardSide{
					Pits: emptySidePits,
				},
				BoardSide2: BoardSide{
					Pits: fullSidePits,
				},
			},
			want: true,
		},
		{
			name: "when none of the sides is empty",
			game: Game{
				BoardSide1: BoardSide{
					Pits: buildPits(2, 2, 2, 2, 2, 2),
				},
				BoardSide2: BoardSide{
					Pits: fullSidePits,
				},
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.game.IsDone())
		})
	}
}

func buildPits(a, b, c, d, e, f int64) pitsArray {
	return pitsArray{a, b, c, d, e, f}
}
