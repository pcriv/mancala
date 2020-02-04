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
	SaveGame(engine.Game) error
	GetGame(string) (*engine.Game, error)
}

// ErrNotFound is when a game cannot be found in the database
var ErrNotFound = errors.New("game not found")

// redisRepo is a redis backed repository
type redisRepo struct {
	db *redis.Client
}

// CreateRepo gets a Repo
func CreateRepo(url string) (Repo, error) {
	options, err := redis.ParseURL(url)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	_, err = client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return &redisRepo{db: client}, nil
}

// SaveGame stores a game on the repo
func (repo *redisRepo) SaveGame(game engine.Game) error {
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
func (repo *redisRepo) GetGame(id string) (*engine.Game, error) {
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
