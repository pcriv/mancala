package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pablocrivella/mancala/internal/core"
	"github.com/redis/go-redis/v9"
)

// GameStore is a redis backed game repo.
type GameStore struct {
	db redis.UniversalClient
}

// NewGameStore creates a new GameStore.
func NewGameStore(client redis.UniversalClient) GameStore {
	return GameStore{db: client}
}

// Save stores a game on redis.
func (r GameStore) Save(ctx context.Context, g core.Game) error {
	json, err := json.Marshal(g)
	if err != nil {
		return err
	}
	err = r.db.Set(ctx, g.ID.String(), string(json), time.Hour*2).Err()
	if err != nil {
		return err
	}
	return nil
}

// Find fetches a game with the given ID from redis.
func (r GameStore) Find(ctx context.Context, id string) (core.Game, error) {
	var g core.Game
	val, err := r.db.Get(ctx, id).Result()
	if err != nil {
		return g, errors.Join(core.ErrGameNotFound, fmt.Errorf("%w: cannot find game with id %v", err, id))
	}
	err = json.Unmarshal([]byte(val), &g)
	if err != nil {
		return g, err
	}
	return g, nil
}
