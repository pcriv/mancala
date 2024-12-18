package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/internal/protomap"
	"github.com/pcriv/mancala/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	googleproto "google.golang.org/protobuf/proto"
)

var _ proto.ServiceServer = serviceServer{}

type (
	requestValidator interface {
		Validate(msg googleproto.Message) error
	}

	serviceServer struct {
		service   mancala.Service
		validator requestValidator
	}
)

func (s serviceServer) CreateGame(ctx context.Context, in *proto.CreateGameRequest) (*proto.CreateGameResponse, error) {
	err := s.validator.Validate(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	g, err := s.service.CreateGame(ctx, in.Player1, in.Player2)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error creating game: %v", err))
	}
	return &proto.CreateGameResponse{
		CreatedGame: protomap.Game(g),
	}, nil
}

func (s serviceServer) FindGame(ctx context.Context, in *proto.FindGameRequest) (*proto.FindGameResponse, error) {
	err := s.validator.Validate(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	g, err := s.service.FindGame(ctx, in.Id)
	if err != nil {
		switch {
		case errors.Is(err, mancala.ErrGameNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("error finding game: %v", err))
		}
	}

	return &proto.FindGameResponse{
		Game: protomap.Game(g),
	}, nil
}

func (s serviceServer) ExecutePlay(ctx context.Context, in *proto.ExecutePlayRequest) (*proto.ExecutePlayResponse, error) {
	err := s.validator.Validate(in)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	g, err := s.service.ExecutePlay(ctx, in.GameId, in.PitIndex)
	if err != nil {
		switch {
		case errors.Is(err, mancala.ErrGameNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, mancala.ErrInvalidPlay):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, fmt.Sprintf("error executing play: %v", err))
		}
	}

	return &proto.ExecutePlayResponse{
		PlayedPitIndex: in.PitIndex,
		Game:           protomap.Game(g),
	}, nil
}
