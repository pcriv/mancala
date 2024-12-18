package mancala

import (
	"context"
)

type gameStore interface {
	Find(ctx context.Context, id string) (Game, error)
	Save(ctx context.Context, g Game) error
}

type Service struct {
	gameStore gameStore
}

func NewService(store gameStore) Service {
	return Service{gameStore: store}
}

func (s Service) CreateGame(ctx context.Context, player1, player2 string) (Game, error) {
	g := NewGame(player1, player2)
	if err := s.gameStore.Save(ctx, g); err != nil {
		return Game{}, err
	}
	return g, nil
}

func (s Service) FindGame(ctx context.Context, id string) (Game, error) {
	g, err := s.gameStore.Find(ctx, id)
	if err != nil {
		return Game{}, err
	}
	return g, nil
}

func (s Service) ExecutePlay(ctx context.Context, gameID string, pitIndex int64) (Game, error) {
	g, err := s.gameStore.Find(ctx, gameID)
	if err != nil {
		return Game{}, err
	}

	err = g.PlayTurn(pitIndex)
	if err != nil {
		return Game{}, err
	}

	err = s.gameStore.Save(ctx, g)
	if err != nil {
		return Game{}, err
	}
	return g, nil
}
