package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pablocrivella/mancala/internal/engine"
	"github.com/pablocrivella/mancala/internal/pkg/dbschema"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type GameStore struct {
	db *sql.DB
}

func NewGameRepo(db *sql.DB) GameStore {
	return GameStore{db: db}
}

// Save stores a game on redis.
func (s GameStore) Save(g engine.Game) error {
	game := &dbschema.Game{
		ID:     g.ID.String(),
		Turn:   int64(g.Turn),
		Result: int64(g.Result),
	}
	game.BoardSide1.Marshal(g.BoardSide1)
	game.BoardSide2.Marshal(g.BoardSide2)
	err := game.Upsert(context.Background(), s.db, true, []string{}, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (s GameStore) Find(id string) (engine.Game, error) {
	g, err := dbschema.FindGame(context.Background(), s.db, id)
	if err != nil {
		return engine.Game{}, err
	}
	var game engine.Game
	game.ID = uuid.MustParse(g.ID)
	game.Result = engine.Result(g.Result)
	game.Turn = engine.Turn(g.Turn)
	g.BoardSide1.Unmarshal(&game.BoardSide1)
	g.BoardSide2.Unmarshal(&game.BoardSide2)
	return game, nil
}
