package protomap

import (
	"github.com/pcriv/mancala/internal/mancala"
	"github.com/pcriv/mancala/proto"
)

func Game(game mancala.Game) *proto.Game {
	return &proto.Game{
		Id:         game.ID.String(),
		BoardSide1: BoardSide(game.BoardSide1),
		BoardSide2: BoardSide(game.BoardSide2),
		Result:     Result(game.Result),
		Turn:       Turn(game.Turn),
	}
}

func Player(player mancala.Player) *proto.Player {
	return &proto.Player{
		Name:  player.Name,
		Score: player.Score,
	}
}

func BoardSide(boardSide mancala.BoardSide) *proto.BoardSide {
	return &proto.BoardSide{
		Pits:   boardSide.Pits[:],
		Store:  boardSide.Store,
		Player: Player(boardSide.Player),
	}
}

func Result(result mancala.Result) proto.Result {
	switch result {
	case mancala.Player1Wins:
		return proto.Result_RESULT_PLAYER1_WINS
	case mancala.Player2Wins:
		return proto.Result_RESULT_PLAYER2_WINS
	case mancala.Tie:
		return proto.Result_RESULT_TIE
	default:
		return proto.Result_RESULT_UNSPECIFIED
	}
}

func Turn(turn mancala.Turn) proto.Turn {
	switch turn {
	case mancala.TurnPlayer1:
		return proto.Turn_TURN_PLAYER1
	case mancala.TurnPlayer2:
		return proto.Turn_TURN_PLAYER2
	default:
		return proto.Turn_TURN_UNSPECIFIED
	}
}
