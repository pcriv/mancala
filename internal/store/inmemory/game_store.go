package inmemory

import (
	"context"

	"github.com/pcriv/mancala/internal/mancala"
)

// GameStore is an im-memory implementation of a GameStore.
type GameStore struct {
	db map[string]mancala.Game
}

func NewGameStore(games ...mancala.Game) GameStore {
	db := make(map[string]mancala.Game)
	for _, g := range games {
		db[g.ID.String()] = g
	}
	return GameStore{db: db}
}

// Find returns a game with the given id.
func (r GameStore) Find(ctx context.Context, id string) (mancala.Game, error) {
	select {
	case <-ctx.Done():
		return mancala.Game{}, ctx.Err()
	default:
	}
	g, ok := r.db[id]
	if ok {
		return g, nil
	}
	return mancala.Game{}, nil
}

// Save persists a game to the in-memory db.
func (r GameStore) Save(ctx context.Context, g mancala.Game) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	r.db[g.ID.String()] = g
	return nil
}
