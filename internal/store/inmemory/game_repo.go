package inmemory

import (
	"context"

	"github.com/pablocrivella/mancala/internal/core"
)

// GameStore is an im-memory implementation of a GameStore.
type GameStore struct {
	db map[string]core.Game
}

func NewGameStore(games ...core.Game) GameStore {
	db := make(map[string]core.Game)
	for _, g := range games {
		db[g.ID.String()] = g
	}
	return GameStore{db: db}
}

// Find returns a game with the given id.
func (r GameStore) Find(ctx context.Context, id string) (core.Game, error) {
	select {
	case <-ctx.Done():
		return core.Game{}, ctx.Err()
	default:
	}
	g, ok := r.db[id]
	if ok {
		return g, nil
	}
	return core.Game{}, nil
}

// Save persists a game to the in-memory db.
func (r GameStore) Save(ctx context.Context, g core.Game) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	r.db[g.ID.String()] = g
	return nil
}
