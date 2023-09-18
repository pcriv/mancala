package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame_PlayTurn(t *testing.T) {
	testCases := []struct {
		name         string
		game         Game
		pitIndex     int64
		wantedGame   Game
		wantedErrMsg string
	}{
		{
			name:     "when player1 captures one stone and has to play another turn",
			pitIndex: 0,
			game:     NewGame("Rick", "Morty"),
			wantedGame: Game{
				BoardSide1: BoardSide{
					Pits:  buildPits(0, 7, 7, 7, 7, 7),
					Store: 1,
				},
				BoardSide2: BoardSide{
					Pits:  fullSidePits,
					Store: 0,
				},
				Turn: Player1Turn,
			},
			wantedErrMsg: "",
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
			},
			wantedGame: Game{
				BoardSide1: BoardSide{
					Pits:  buildPits(0, 0, 0, 0, 0, 1),
					Store: 2,
				},
				BoardSide2: BoardSide{
					Pits:  buildPits(6, 6, 6, 6, 0, 6),
					Store: 0,
				},
				Turn: Player2Turn,
			},
			wantedErrMsg: "",
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
			},
			wantedGame: Game{
				BoardSide1: BoardSide{
					Pits:  buildPits(2, 1, 0, 0, 0, 0),
					Store: 1,
				},
				BoardSide2: BoardSide{
					Pits:  buildPits(7, 7, 7, 7, 7, 7),
					Store: 0,
				},
				Turn: Player2Turn,
			},
			wantedErrMsg: "",
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
			},
			wantedGame: Game{
				BoardSide1: BoardSide{
					Pits:  buildPits(0, 0, 0, 0, 0, 0),
					Store: 1,
				},
				BoardSide2: BoardSide{
					Pits:  buildPits(1, 1, 1, 1, 1, 1),
					Store: 6,
				},
				Turn:   Player1Turn,
				Result: Player2Wins,
			},
			wantedErrMsg: "",
		},
		{
			name:         "when the pit index is invalid",
			pitIndex:     8,
			game:         Game{},
			wantedGame:   Game{},
			wantedErrMsg: "invalid play: pit index is invalid",
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
			},
			wantedGame:   Game{},
			wantedErrMsg: "invalid play: selected pit is empty",
		},
		{
			name:     "the game is already done",
			pitIndex: 0,
			game: Game{
				BoardSide1: BoardSide{
					Pits: buildPits(0, 0, 0, 0, 0, 1),
				},
			},
			wantedGame:   Game{},
			wantedErrMsg: "invalid play: game is already done",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.game.PlayTurn(tc.pitIndex)
			if tc.wantedErrMsg == "" {
				assert.Equal(t, tc.wantedGame.BoardSide1.Pits, tc.game.BoardSide1.Pits)
				assert.Equal(t, tc.wantedGame.BoardSide2.Pits, tc.game.BoardSide2.Pits)
				assert.Equal(t, tc.wantedGame.BoardSide1.Store, tc.game.BoardSide1.Store)
				assert.Equal(t, tc.wantedGame.BoardSide2.Store, tc.game.BoardSide2.Store)
				assert.Equal(t, tc.wantedGame.Turn, tc.game.Turn)
				assert.Equal(t, tc.wantedGame.Result, tc.game.Result)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tc.wantedErrMsg)
				}
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
