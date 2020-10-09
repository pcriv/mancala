package main

import (
	"net/http"
	"os"

	"github.com/pablocrivella/mancala/api/rpc"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
	"github.com/pablocrivella/mancala/internal/rpcapi"
)

func main() {
	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		panic("missing env variable: REDIS_URL")
	}
	redisClient, err := persistence.NewRedisClient(redisURL)
	if err != nil {
		panic(err)
	}
	gameRepo := persistence.NewGameRepo(redisClient)
	gamesService := games.NewService(gameRepo)
	server := &rpcapi.Server{GamesService: gamesService}
	twirpHandler := rpc.NewMancalaServer(server)

	http.ListenAndServe(":8080", twirpHandler)
}
