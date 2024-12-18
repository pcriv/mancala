package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pcriv/mancala/internal/mancala"

	"github.com/alicebob/miniredis/v2"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	redisstore "github.com/pcriv/mancala/internal/store/redis"
)

func TestGamesHandler_Create(t *testing.T) {
	s, closeRedisFunc := startFakeRedisServer(t)
	defer closeRedisFunc()

	redisClient := newRedisClient(t, "redis://"+s.Addr())
	gameStore := redisstore.NewGameStore(redisClient)
	gameService := mancala.NewService(gameStore)
	h := GamesHandler{GameService: gameService}
	e := echo.New()

	testCases := []struct {
		name       string
		body       string
		wantedCode int
	}{
		{
			name:       "when the request is valid",
			body:       `{"player1": "Rick","player2": "Morty"}`,
			wantedCode: 201,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rq := httptest.NewRequest(http.MethodPost, "/v1/games", strings.NewReader(tc.body))
			rq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rc := httptest.NewRecorder()

			if assert.NoError(t, h.CreateGame(e.NewContext(rq, rc))) {
				assert.Equal(t, tc.wantedCode, rc.Code)
				assert.NotEmpty(t, strings.TrimSpace(rc.Body.String()))
			}
		})
	}
}

func TestGamesHandler_Show(t *testing.T) {
	s, closeRedisFunc := startFakeRedisServer(t)
	defer closeRedisFunc()

	redisClient := newRedisClient(t, "redis://"+s.Addr())
	gameStore := redisstore.NewGameStore(redisClient)
	gameService := mancala.NewService(gameStore)
	h := GamesHandler{GameService: gameService}

	g, err := gameService.CreateGame(context.Background(), "Rick", "Morty")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	e := echo.New()

	testCases := []struct {
		name       string
		game       mancala.Game
		body       string
		wantedCode int
	}{
		{
			name:       "when the game exists",
			game:       g,
			wantedCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rq := httptest.NewRequest(http.MethodGet, "/v1/games/:id", nil)
			rq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rc := httptest.NewRecorder()
			ctx := e.NewContext(rq, rc)
			ctx.SetPath("/v1/games/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(tc.game.ID.String())

			if assert.NoError(t, h.ShowGame(ctx, tc.game.ID.String())) {
				assert.Equal(t, tc.wantedCode, rc.Code)
				assert.NotEmpty(t, strings.TrimSpace(rc.Body.String()))
			}
		})
	}
}

func startFakeRedisServer(t *testing.T) (*miniredis.Miniredis, func()) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	return s, func() { s.Close() }
}

func newRedisClient(t *testing.T, url string) *redis.Client {
	options, err := redis.ParseURL(url)
	require.NoError(t, err)

	client := redis.NewClient(options)
	_, err = client.Ping(context.Background()).Result()
	require.NoError(t, err)

	return client
}
