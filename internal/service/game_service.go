package service

import (
	"context"

	"github.com/pablocrivella/mancala/internal/core"
)

// Service is a games context domain service
type GameService struct {
	gameStore core.GameStore
}

func NewGameService(store core.GameStore) GameService {
	return GameService{gameStore: store}
}

func (s GameService) CreateGame(ctx context.Context, player1, player2 string) (core.Game, error) {
	g := core.NewGame(player1, player2)
	if err := s.gameStore.Save(ctx, g); err != nil {
		return core.Game{}, err
	}
	return g, nil
}

func (s GameService) FindGame(ctx context.Context, id string) (core.Game, error) {
	g, err := s.gameStore.Find(ctx, id)
	if err != nil {
		return core.Game{}, err
	}
	return g, nil
}

func (s GameService) ExecutePlay(ctx context.Context, gameID string, pitIndex int64) (core.Game, error) {
	g, err := s.gameStore.Find(ctx, gameID)
	if err != nil {
		return core.Game{}, err
	}

	err = g.PlayTurn(pitIndex)
	if err != nil {
		return core.Game{}, err
	}

	err = s.gameStore.Save(ctx, g)
	if err != nil {
		return core.Game{}, err
	}
	return g, nil
}
