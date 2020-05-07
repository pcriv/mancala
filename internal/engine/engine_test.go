package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	game := NewGame("Rick", "Morty")

	assert.Equal(t, game.BoardSide1.Pits, [6]int{6, 6, 6, 6, 6, 6})
	assert.Equal(t, game.BoardSide2.Pits, [6]int{6, 6, 6, 6, 6, 6})

	assert.Equal(t, game.BoardSide1.Player.Name, "Rick")
	assert.Equal(t, game.BoardSide2.Player.Name, "Morty")
}
