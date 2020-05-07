package engine

type (
	// Player represents a player of the mancala game
	Player struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}
)
