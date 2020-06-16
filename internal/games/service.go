package games

import (
	"github.com/pablocrivella/mancala/internal/engine"
)

// Service is a games context domain service
type Service struct {
	gameRepo GameRepo
}

// NewService returns a games.Service
func NewService(g GameRepo) Service {
	return Service{gameRepo: g}
}

func (s Service) CreateGame(player1, player2 string) (engine.Game, error) {
	g := engine.NewGame(player1, player2)
	if err := s.gameRepo.Save(g); err != nil {
		return engine.Game{}, err
	}
	return g, nil
}

func (s Service) FindGame(id string) (engine.Game, error) {
	g, err := s.gameRepo.Find(id)
	if err != nil {
		return engine.Game{}, err
	}
	return g, nil
}

func (s Service) ExecutePlay(gameID string, pitIndex int) (engine.Game, error) {
	g, err := s.gameRepo.Find(gameID)
	if err != nil {
		return engine.Game{}, err
	}

	err = g.PlayTurn(pitIndex)
	if err != nil {
		return engine.Game{}, err
	}

	err = s.gameRepo.Save(g)
	if err != nil {
		return engine.Game{}, err
	}
	return g, nil
}
