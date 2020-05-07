package engine

type (
	// BoardSide represents one of the sides of the board
	BoardSide struct {
		Pits   [pitsPerSide]int `json:"pits"`
		Store  int              `json:"store"`
		Player Player           `json:"player"`
	}
)

func (side *BoardSide) setup(playerName string) {
	side.Player.Name = playerName
	side.Pits = fullSidePits
}

func (side *BoardSide) isEmpty() bool {
	return side.Pits == emptySidePits
}

func (side *BoardSide) score() {
	score := 0

	for _, pitStones := range side.Pits {
		score += pitStones
	}

	side.Store += score
	side.Player.Score = score
}

func (side *BoardSide) pickStones(pitIndex int) int {
	stones := side.Pits[pitIndex]
	side.Pits[pitIndex] = 0

	return stones
}

func (side *BoardSide) placeStone(pitIndex int, stones int) int {
	side.Pits[pitIndex]++

	return stones - 1
}

func (side *BoardSide) captureStone(stones int) int {
	side.Store++

	return stones - 1
}

func (side *BoardSide) captureAllStones(stones int) {
	side.Store += stones
}
