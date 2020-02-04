package persistence

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/pablocrivella/mancala/engine"
)

// Repo interface
type Repo interface {
	Connect(string) error
	SaveGame(engine.Game) error
	GetGame(string) (*engine.Game, error)
}

// ErrNotFound is when a game cannot be found in the database
var ErrNotFound = errors.New("game not found")

// RedisRepo is a redis backed repository
type RedisRepo struct {
	db *redis.Client
}

// Connect to the database
func (repo *RedisRepo) Connect(url string) error {
	options, err := redis.ParseURL(url)

	if err != nil {
		return errors.New(err.Error())
	}

	repo.db = redis.NewClient(options)

	return nil
}

// SaveGame stores a game on the repo
func (repo *RedisRepo) SaveGame(game engine.Game) error {
	json, err := json.Marshal(game)

	if err != nil {
		return err
	}

	err = repo.db.Set(game.ID.String(), string(json), time.Hour*2).Err()

	if err != nil {
		return err
	}

	return nil
}

// GetGame fetches a game from the repo
func (repo *RedisRepo) GetGame(id string) (*engine.Game, error) {
	var game engine.Game

	val, err := repo.db.Get(id).Result()

	if err != nil {
		return nil, ErrNotFound
	}

	err = json.Unmarshal([]byte(val), &game)

	if err != nil {
		return nil, err
	}

	return &game, nil
}
