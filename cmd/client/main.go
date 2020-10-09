package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pablocrivella/mancala/api/rpc"
)

func main() {
	client := rpc.NewMancalaProtobufClient("http://localhost:8080", &http.Client{})
	game, err := client.CreateGame(context.Background(), &rpc.NewGame{Player1: "Red", Player2: "Blue"})
	if err != nil {
		fmt.Printf("oh no: %v", err)
		os.Exit(1)
	}
	fmt.Printf("New game created: %+v", game)
}

// curl --request "POST" --header "Content-Type: application/json" --data '{"gameId": "placeholder",  "pitIndex": 0}' http://localhost:8080/twirp/rpc.Mancala/ExecutePlay
// curl --request "POST" --header "Content-Type: application/json" --data '{"player1": "Red",  "player2": "Blue"}' http://localhost:8080/twirp/rpc.Mancala/CreateGame
// curl --request "POST" --header "Content-Type: application/json" --data '{"id": "placeholder"}' http://localhost:8080/twirp/rpc.Mancala/GetGame
