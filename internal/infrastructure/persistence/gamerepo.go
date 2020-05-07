package persistence

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/pablocrivella/mancala/internal/engine"
)

type (
	// GameRepo is a redis backed game repo.
	GameRepo struct {
		db *redis.Client
	}
)

// NewGameRepo creates a new GameRepo.
func NewGameRepo(url string) (*GameRepo, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &GameRepo{db: client}, nil
}

// Save stores a game on redis.
func (r *GameRepo) Save(g *engine.Game) error {
	json, err := json.Marshal(g)
	if err != nil {
		return err
	}

	err = r.db.Set(g.ID.String(), string(json), time.Hour*2).Err()
	if err != nil {
		return err
	}
	return nil
}

// Find fetches a game with the given ID from redis.
func (r *GameRepo) Find(id string) (*engine.Game, error) {
	var g engine.Game
	val, err := r.db.Get(id).Result()
	if err != nil {
		return nil, &ErrNotFound{Msg: fmt.Sprintf("cannot find with id %v", id)}
	}

	err = json.Unmarshal([]byte(val), &g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}
