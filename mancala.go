package mancala

// NewGame initializes a manala game for the given players
func NewGame(player1Name string, player2Name string) Game {
	game := Game{}
	game.Setup(player1Name, player2Name)

	return game
}
