package rpcapi

import (
	"context"
	"strings"

	"github.com/pablocrivella/mancala/api/rpc"
	"github.com/pablocrivella/mancala/internal/engine"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/twitchtv/twirp"
)

type Server struct {
	GamesService games.Service
}

func (s *Server) CreateGame(ctx context.Context, newGame *rpc.NewGame) (*rpc.OngoingGame, error) {
	if strings.TrimSpace(newGame.Player1) == "" {
		return nil, twirp.RequiredArgumentError("player1")
	}
	if strings.TrimSpace(newGame.Player2) == "" {
		return nil, twirp.RequiredArgumentError("player2")
	}
	g, err := s.GamesService.CreateGame(newGame.Player1, newGame.Player2)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}
	return covertToOngoingGameMessage(g), nil
}

func (s *Server) ExecutePlay(ctx context.Context, p *rpc.Play) (*rpc.OngoingGame, error) {
	if strings.TrimSpace(p.GameId) == "" {
		return nil, twirp.RequiredArgumentError("gameId")
	}
	g, err := s.GamesService.ExecutePlay(p.GameId, p.PitIndex)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}
	return covertToOngoingGameMessage(g), nil
}

func (s *Server) GetGame(ctx context.Context, r *rpc.GameRequest) (*rpc.OngoingGame, error) {
	if strings.TrimSpace(r.Id) == "" {
		return nil, twirp.RequiredArgumentError("id")
	}
	g, err := s.GamesService.FindGame(r.Id)
	if err != nil {
		return nil, twirp.InternalErrorWith(err)
	}
	return covertToOngoingGameMessage(g), nil
}

func covertToOngoingGameMessage(g engine.Game) *rpc.OngoingGame {
	return &rpc.OngoingGame{
		Id:     g.ID.String(),
		Result: int64(g.Result),
		Turn:   int64(g.Turn),
		Board1: &rpc.BoardState{
			Pits:  g.BoardSide1.Pits[:],
			Store: g.BoardSide1.Store,
			Player: &rpc.Player{
				Name:  g.BoardSide1.Player.Name,
				Score: g.BoardSide1.Player.Score,
			},
		},
		Board2: &rpc.BoardState{
			Pits:  g.BoardSide2.Pits[:],
			Store: g.BoardSide2.Store,
			Player: &rpc.Player{
				Name:  g.BoardSide1.Player.Name,
				Score: g.BoardSide1.Player.Score,
			},
		},
	}
}
