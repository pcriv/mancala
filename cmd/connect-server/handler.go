package main

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/internal/protomap"
	"github.com/pcriv/mancala/proto"
	"github.com/pcriv/mancala/proto/protoconnect"
	googleproto "google.golang.org/protobuf/proto"
)

var _ protoconnect.ServiceHandler = handler{}

type (
	requestValidator interface {
		Validate(msg googleproto.Message) error
	}

	handler struct {
		service   mancala.Service
		validator requestValidator
	}
)

func (h handler) CreateGame(ctx context.Context, in *connect.Request[proto.CreateGameRequest]) (*connect.Response[proto.CreateGameResponse], error) {
	err := h.validator.Validate(in.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	g, err := h.service.CreateGame(ctx, in.Msg.Player1, in.Msg.Player2)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&proto.CreateGameResponse{
		CreatedGame: protomap.Game(g),
	}), nil
}

func (h handler) FindGame(ctx context.Context, in *connect.Request[proto.FindGameRequest]) (*connect.Response[proto.FindGameResponse], error) {
	err := h.validator.Validate(in.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	g, err := h.service.FindGame(ctx, in.Msg.Id)
	if err != nil {
		switch {
		case errors.Is(err, mancala.ErrGameNotFound):
			return nil, connect.NewError(connect.CodeNotFound, err)
		default:
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&proto.FindGameResponse{
		Game: protomap.Game(g),
	}), nil
}

func (h handler) ExecutePlay(ctx context.Context, in *connect.Request[proto.ExecutePlayRequest]) (*connect.Response[proto.ExecutePlayResponse], error) {
	err := h.validator.Validate(in.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	g, err := h.service.ExecutePlay(ctx, in.Msg.GameId, in.Msg.PitIndex)
	if err != nil {
		switch {
		case errors.Is(err, mancala.ErrGameNotFound):
			return nil, connect.NewError(connect.CodeNotFound, err)
		case errors.Is(err, mancala.ErrInvalidPlay):
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		default:
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&proto.ExecutePlayResponse{
		PlayedPitIndex: in.Msg.PitIndex,
		Game:           protomap.Game(g),
	}), nil
}
