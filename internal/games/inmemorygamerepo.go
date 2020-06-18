package games

import (
	"github.com/pablocrivella/mancala/internal/engine"
)

// InMemoryGameRepo is an im-memory implementation of a GameRepo.
type InMemoryGameRepo struct {
	db map[string]engine.Game
}

// NewFakeGameRepo returns an in-memory implementation of a GameRepo.
func NewFakeGameRepo(games ...engine.Game) GameRepo {
	db := make(map[string]engine.Game)
	for _, g := range games {
		db[g.ID.String()] = g
	}
	return InMemoryGameRepo{db: db}
}

// Find returns a game with the given id.
func (r InMemoryGameRepo) Find(id string) (engine.Game, error) {
	g, ok := r.db[id]
	if ok {
		return g, nil
	}
	return engine.Game{}, nil
}

// Save persists a game to the in-memory db.
func (r InMemoryGameRepo) Save(g engine.Game) error {
	r.db[g.ID.String()] = g
	return nil
}
