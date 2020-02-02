package repo

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/pablocrivella/mancala/engine"
	"golang.org/x/xerrors"
)

var db *redis.Client

// ErrNotFound is when a game cannot be found in the database
var ErrNotFound = errors.New("game not found")

// Connect to the database
func Connect() error {
	url, ok := os.LookupEnv("REDIS_URL")

	if !ok {
		return xerrors.New("missing env variable: REDIS_URL")
	}

	options, err := redis.ParseURL(url)

	if err != nil {
		return xerrors.New(err.Error())
	}

	db = redis.NewClient(options)

	return nil
}

// SaveGame stores a game on the repo
func SaveGame(game engine.Game) error {
	json, err := json.Marshal(game)

	if err != nil {
		return err
	}

	err = Connect()

	if err != nil {
		return err
	}

	err = db.Set(game.ID.String(), string(json), time.Hour*2).Err()

	if err != nil {
		return err
	}

	return nil
}

// GetGame fetches a game from the repo
func GetGame(id string) (*engine.Game, error) {
	var game engine.Game

	err := Connect()

	if err != nil {
		return nil, err
	}

	val, err := db.Get(id).Result()

	if err != nil {
		return nil, ErrNotFound
	}

	err = json.Unmarshal([]byte(val), &game)

	if err != nil {
		return nil, err
	}

	return &game, nil
}
