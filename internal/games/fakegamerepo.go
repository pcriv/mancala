package games

import (
	"github.com/pablocrivella/mancala/internal/engine"
)

type fakeGameRepo struct {
	store map[string]engine.Game
}

// NewFakeGameRepo returns a fake in-memory implementation of a GameRepo.
func NewFakeGameRepo(games ...engine.Game) GameRepo {
	store := make(map[string]engine.Game)
	for _, g := range games {
		store[g.ID.String()] = g
	}
	return fakeGameRepo{store: store}
}

// Find returns a game with the given id.
func (r fakeGameRepo) Find(id string) (engine.Game, error) {
	g, ok := r.store[id]
	if ok {
		return g, nil
	}
	return engine.Game{}, nil
}

// Save persists a game to the in-memory store.
func (r fakeGameRepo) Save(g engine.Game) error {
	r.store[g.ID.String()] = g
	return nil
}
