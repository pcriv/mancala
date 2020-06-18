package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pablocrivella/mancala/internal/engine"
)

type (
	// GameRepo is a redis backed game repo.
	GameRepo struct {
		db *redis.Client
	}
)

// NewGameRepo creates a new GameRepo.
func NewGameRepo(client *redis.Client) GameRepo {
	return GameRepo{db: client}
}

// Save stores a game on redis.
func (r GameRepo) Save(g engine.Game) error {
	json, err := json.Marshal(g)
	if err != nil {
		return err
	}
	err = r.db.Set(context.Background(), g.ID.String(), string(json), time.Hour*2).Err()
	if err != nil {
		return err
	}
	return nil
}

// Find fetches a game with the given ID from redis.
func (r GameRepo) Find(id string) (engine.Game, error) {
	var g engine.Game
	val, err := r.db.Get(context.Background(), id).Result()
	if err != nil {
		return g, &NotFoundError{Msg: fmt.Sprintf("cannot find game with id %v", id)}
	}
	err = json.Unmarshal([]byte(val), &g)
	if err != nil {
		return g, err
	}
	return g, nil
}
