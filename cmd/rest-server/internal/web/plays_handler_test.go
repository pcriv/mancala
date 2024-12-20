package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pcriv/mancala/internal/mancala"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	redisstore "github.com/pcriv/mancala/internal/store/redis"
)

func TestPlaysHandler_Create(t *testing.T) {
	s, closeRedisFunc := startFakeRedisServer(t)
	defer closeRedisFunc()

	redisClient := newRedisClient(t, "redis://"+s.Addr())
	gameStore := redisstore.NewGameStore(redisClient)
	gameService := mancala.NewService(gameStore)
	h := PlaysHandler{GameService: gameService}
	e := echo.New()
	g, err := gameService.CreateGame(context.Background(), "Rick", "Morty")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCases := []struct {
		name       string
		game       mancala.Game
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
			rq := httptest.NewRequest(http.MethodPost, "/v1/games/:id/plays", strings.NewReader(tc.body))
			rq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rc := httptest.NewRecorder()
			ctx := e.NewContext(rq, rc)
			ctx.SetPath("/v1/games/:id/plays")
			ctx.SetParamNames("id")
			ctx.SetParamValues(tc.game.ID.String())

			if assert.NoError(t, h.CreatePlay(ctx, tc.game.ID.String())) {
				assert.Equal(t, tc.wantedCode, rc.Code)
				assert.NotEmpty(t, strings.TrimSpace(rc.Body.String()))
			}
		})
	}
}
