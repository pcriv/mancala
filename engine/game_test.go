package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayTurn(t *testing.T) {
	game := NewGame("Rick", "Morty")

	// Tests scenario where Player1 captures one stone and has to play another turn.
	game.PlayTurn(0)

	assert.Equal(t, game.BoardSide1.Pits, [6]int{0, 7, 7, 7, 7, 7})
	assert.Equal(t, game.BoardSide2.Pits, [6]int{6, 6, 6, 6, 6, 6})

	assert.Equal(t, game.BoardSide1.Store, 1)
	assert.Equal(t, game.BoardSide2.Store, 0)
	assert.Equal(t, game.Turn, Player1Turn)

	// Test scenario where Player1 moves the last stone to an empty own pit.

	game = NewGame("Rick", "Morty")
	game.BoardSide1.Pits = [6]int{1, 0, 0, 0, 0, 1}
	game.BoardSide2.Pits = [6]int{6, 6, 6, 6, 1, 6}
	game.PlayTurn(0)

	assert.Equal(t, game.BoardSide1.Pits, [6]int{0, 0, 0, 0, 0, 1})
	assert.Equal(t, game.BoardSide2.Pits, [6]int{6, 6, 6, 6, 0, 6})

	assert.Equal(t, game.BoardSide1.Store, 2)
	assert.Equal(t, game.BoardSide2.Store, 0)
	assert.Equal(t, game.Turn, Player2Turn)

	// Test scenario where Player1 places stones on Player2 pit but skipping Player2 store.

	game = NewGame("Rick", "Morty")
	game.BoardSide1.Pits = [6]int{1, 1, 0, 0, 0, 8}
	game.BoardSide2.Pits = [6]int{6, 6, 6, 6, 6, 6}
	game.PlayTurn(5)

	assert.Equal(t, game.BoardSide1.Pits, [6]int{2, 1, 0, 0, 0, 0})
	assert.Equal(t, game.BoardSide2.Pits, [6]int{7, 7, 7, 7, 7, 7})

	assert.Equal(t, game.BoardSide1.Store, 1)
	assert.Equal(t, game.BoardSide2.Store, 0)
	assert.Equal(t, game.Turn, Player2Turn)

	// Test scenario where Player2 captures a stone.

	game = NewGame("Rick", "Morty")
	game.BoardSide1.Pits = [6]int{2, 1, 0, 0, 0, 0}
	game.BoardSide2.Pits = [6]int{7, 7, 7, 7, 7, 7}
	game.Turn = Player2Turn

	game.PlayTurn(5)

	assert.Equal(t, game.BoardSide1.Pits, [6]int{3, 2, 1, 1, 1, 1})
	assert.Equal(t, game.BoardSide2.Pits, [6]int{7, 7, 7, 7, 7, 0})

	assert.Equal(t, game.BoardSide1.Store, 0)
	assert.Equal(t, game.BoardSide2.Store, 1)
	assert.Equal(t, game.Turn, Player1Turn)

	// Test scenario where it is the last Turn before the end of the game

	game = NewGame("Rick", "Morty")
	game.BoardSide1.Pits = [6]int{0, 0, 0, 0, 0, 1}
	game.BoardSide2.Pits = [6]int{1, 1, 1, 1, 1, 1}

	game.PlayTurn(5)

	assert.Equal(t, game.BoardSide1.Pits, [6]int{0, 0, 0, 0, 0, 0})
	assert.Equal(t, game.BoardSide2.Pits, [6]int{1, 1, 1, 1, 1, 1})

	assert.Equal(t, game.BoardSide1.Store, 1)
	assert.Equal(t, game.BoardSide2.Store, 6)
	assert.Equal(t, game.Result, Player2Wins)

	// Test scenario where the pitIndex is invalid

	game = NewGame("Rick", "Morty")
	err := game.PlayTurn(8)

	assert.Equal(t, err.Error(), "pit index is invalid")

	// Test scenario where the selected pit is empty

	game = NewGame("Rick", "Morty")
	game.BoardSide1.Pits = [6]int{0, 0, 0, 0, 0, 1}
	err = game.PlayTurn(0)

	assert.Equal(t, err.Error(), "selected pit is empty")

	// Test scenario where the game is already done

	game = NewGame("Rick", "Morty")
	game.BoardSide1.Pits = [6]int{0, 0, 0, 0, 0, 1}
	game.BoardSide2.Pits = [6]int{1, 1, 1, 1, 1, 1}

	game.PlayTurn(5)

	err = game.PlayTurn(0)

	assert.Equal(t, err.Error(), "game is already done")
}

func TestIsDone(t *testing.T) {
	game := NewGame("Rick", "Morty")

	// When one of the sides is empty
	game.BoardSide1.Pits = [6]int{}
	game.BoardSide2.Pits = [6]int{1, 1, 1, 1, 1, 1}

	assert.True(t, game.IsDone())

	// When none of the sides is empty
	game.BoardSide1.Pits = [6]int{2, 2, 2, 2, 2, 2}
	game.BoardSide2.Pits = [6]int{1, 1, 1, 1, 1, 1}

	assert.False(t, game.IsDone())
}
