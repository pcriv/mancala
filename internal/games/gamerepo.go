package games

import "github.com/pablocrivella/mancala/internal/engine"

type GameRepo interface {
	Find(string) (engine.Game, error)
	Save(engine.Game) error
}
