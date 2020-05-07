package resources

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/internal/engine"
	"github.com/pablocrivella/mancala/internal/games"
	"github.com/pablocrivella/mancala/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
)

var (
	createGameReqJSON = `
		{
			"player1": "Rick",
			"player2": "Morty"
		}
	`
	updateGameReqJSON = `
		{
			"pit_index": 0
		}
	`
)

func TestGamesResource_Create(t *testing.T) {
	// Setup
	e := echo.New()
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	gr, err := persistence.NewGameRepo("redis://" + s.Addr())
	if err != nil {
		panic(err)
	}

	h := GamesResource{GamesService: games.NewService(gr)}
	req := httptest.NewRequest(http.MethodPost, "/v1/games", strings.NewReader(createGameReqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, h.Create(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var game engine.Game

		json.Unmarshal(rec.Body.Bytes(), &game)

		assert.NotNil(t, game.ID)
		assert.Equal(t, "Rick", game.BoardSide1.Player.Name)
		assert.Equal(t, "Morty", game.BoardSide2.Player.Name)
	}
}

func TestGamesResource_Update(t *testing.T) {
	// Setup
	e := echo.New()
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	gr, err := persistence.NewGameRepo("redis://" + s.Addr())
	if err != nil {
		panic(err)
	}

	gs := games.NewService(gr)
	g, err := gs.CreateGame("Rick", "Morty")
	if err != nil {
		panic(err)
	}

	h := GamesResource{GamesService: gs}
	req := httptest.NewRequest(http.MethodPatch, "/v1/games/:id", strings.NewReader(updateGameReqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/games/:id")
	c.SetParamNames("id")
	c.SetParamValues(g.ID.String())

	// Assertions
	if assert.NoError(t, h.Update(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var game engine.Game

		json.Unmarshal(rec.Body.Bytes(), &game)

		assert.NotNil(t, g.ID, game.ID)
		assert.Equal(t, engine.Player1Turn, game.Turn)
		assert.Equal(t, engine.Undefined, game.Result)
		assert.Equal(t, [6]int{0, 7, 7, 7, 7, 7}, game.BoardSide1.Pits)
		assert.Equal(t, [6]int{6, 6, 6, 6, 6, 6}, game.BoardSide2.Pits)
		assert.Equal(t, 1, game.BoardSide1.Store)
		assert.Equal(t, 0, game.BoardSide2.Store)
		assert.Equal(t, "Rick", game.BoardSide1.Player.Name)
		assert.Equal(t, "Morty", game.BoardSide2.Player.Name)
	}
}

func TestGamesResorce_Show(t *testing.T) {
	// Setup
	e := echo.New()
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	gr, err := persistence.NewGameRepo("redis://" + s.Addr())
	if err != nil {
		panic(err)
	}

	gs := games.NewService(gr)
	g, err := gs.CreateGame("Rick", "Morty")
	if err != nil {
		panic(err)
	}

	h := GamesResource{GamesService: gs}
	req := httptest.NewRequest(http.MethodGet, "/v1/games/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/games/:id")
	c.SetParamNames("id")
	c.SetParamValues(g.ID.String())

	// Assertions
	if assert.NoError(t, h.Show(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var game engine.Game

		json.Unmarshal(rec.Body.Bytes(), &game)

		assert.Equal(t, g.ID, game.ID)
		assert.Equal(t, "Rick", game.BoardSide1.Player.Name)
		assert.Equal(t, "Morty", game.BoardSide2.Player.Name)
	}
}
