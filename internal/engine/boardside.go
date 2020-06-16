package engine

type (
	pitsArray [pitsPerSide]int
	// BoardSide represents one of the sides of the board
	BoardSide struct {
		Pits   pitsArray `json:"pits"`
		Store  int       `json:"store"`
		Player Player    `json:"player"`
	}
)

var fullSidePits = pitsArray{6, 6, 6, 6, 6, 6}
var emptySidePits = pitsArray{}

func (side *BoardSide) setup(playerName string) {
	side.Player.Name = playerName
	side.Pits = fullSidePits
}

func (side *BoardSide) isEmpty() bool {
	for _, pit := range side.Pits {
		if pit != 0 {
			return false
		}
	}
	return true
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
