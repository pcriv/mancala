package resources

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/engine"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
)

func TestGamesResource_Create(t *testing.T) {
	s, closeRedisFunc := startFakeRedisServer(t)
	defer closeRedisFunc()

	redisClient := newRedisClient(t, "redis://"+s.Addr())
	gameRepo := persistence.NewGameRepo(redisClient)
	h := GamesResource{GamesService: games.NewService(gameRepo)}
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

func TestGamesResource_Update(t *testing.T) {
	s, closeRedisFunc := startFakeRedisServer(t)
	defer closeRedisFunc()

	redisClient := newRedisClient(t, "redis://"+s.Addr())
	gameRepo := persistence.NewGameRepo(redisClient)
	gamesService := games.NewService(gameRepo)
	h := GamesResource{GamesService: gamesService}
	e := echo.New()
	g, err := gamesService.CreateGame("Rick", "Morty")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCases := []struct {
		name       string
		game       engine.Game
		body       string
		wantedCode int
	}{
		{
			name:       "",
			game:       g,
			body:       `{"pit_index": 0}`,
			wantedCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rq := httptest.NewRequest(http.MethodPut, "/v1/games/:id", strings.NewReader(tc.body))
			rq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rc := httptest.NewRecorder()
			ctx := e.NewContext(rq, rc)
			ctx.SetPath("/v1/games/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues(tc.game.ID.String())

			if assert.NoError(t, h.Update(ctx, tc.game.ID.String())) {
				assert.Equal(t, tc.wantedCode, rc.Code)
				assert.NotEmpty(t, strings.TrimSpace(rc.Body.String()))
			}
		})
	}
}

func TestGamesResorce_Show(t *testing.T) {
	s, closeRedisFunc := startFakeRedisServer(t)
	defer closeRedisFunc()

	redisClient := newRedisClient(t, "redis://"+s.Addr())
	gameRepo := persistence.NewGameRepo(redisClient)
	gamesService := games.NewService(gameRepo)
	h := GamesResource{GamesService: gamesService}
	g, err := gamesService.CreateGame("Rick", "Morty")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	e := echo.New()

	testCases := []struct {
		name      string
		game      engine.Game
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
	c, err := persistence.NewRedisClient(url)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return c
}
