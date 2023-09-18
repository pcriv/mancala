package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	game := NewGame("Rick", "Morty")

	assert.Equal(t, game.BoardSide1.Pits, fullSidePits)
	assert.Equal(t, game.BoardSide2.Pits, fullSidePits)
	assert.Equal(t, game.BoardSide1.Player.Name, "Rick")
	assert.Equal(t, game.BoardSide2.Player.Name, "Morty")
	assert.Equal(t, game.Turn, Player1Turn)
	assert.Equal(t, game.Result, Undefined)
}
