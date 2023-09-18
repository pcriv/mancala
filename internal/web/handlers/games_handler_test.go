package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/core"
	"github.com/pablocrivella/mancala/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisstore "github.com/pablocrivella/mancala/internal/store/redis"
)

func TestGamesHandler_Create(t *testing.T) {
	s, closeRedisFunc := startFakeRedisServer(t)
	defer closeRedisFunc()

	redisClient := newRedisClient(t, "redis://"+s.Addr())
	gameStore := redisstore.NewGameStore(redisClient)
	gameService := service.NewGameService(gameStore)
	h := GamesHandler{GameService: gameService}
	e := echo.New()

	testCases := []struct {
		name      string
		body      string
		watedCode int
	}{
		{
			name:      "when the request is valid",
			body:      `{"player1": "Rick","player2": "Morty"}`,
			watedCode: 201,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rq := httptest.NewRequest(http.MethodPost, "/v1/games", strings.NewReader(tc.body))
			rq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rc := httptest.NewRecorder()

			if assert.NoError(t, h.Create(e.NewContext(rq, rc))) {
				assert.Equal(t, tc.watedCode, rc.Code)
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
	gameService := service.NewGameService(gameStore)
	h := GamesHandler{GameService: gameService}

	g, err := gameService.CreateGame(context.Background(), "Rick", "Morty")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	e := echo.New()

	testCases := []struct {
		name      string
		game      core.Game
		body      string
		watedCode int
	}{
		{
			name:      "when the game exists",
			game:      g,
			watedCode: 200,
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

			if assert.NoError(t, h.Show(ctx, tc.game.ID.String())) {
				assert.Equal(t, tc.watedCode, rc.Code)
				assert.NotEmpty(t, strings.TrimSpace(rc.Body.String()))
			}
		})
	}
}

func startFakeRedisServer(t *testing.T) (*miniredis.Miniredis, func()) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return s, func() { s.Close() }
}

func newRedisClient(t *testing.T, url string) *redis.Client {
	options, err := redis.ParseURL(url)
	if err != nil {
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}
	client := redis.NewClient(options)
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}
	return client
}
