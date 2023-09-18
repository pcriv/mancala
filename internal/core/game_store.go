package core

import "context"

type GameStore interface {
	Find(context.Context, string) (Game, error)
	Save(context.Context, Game) error
}
