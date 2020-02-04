package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/labstack/echo/v4"
	"github.com/pablocrivella/mancala/engine"
	"github.com/pablocrivella/mancala/persistence"
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

func TestCreateGame(t *testing.T) {
	// Setup
	e := echo.New()

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	repo := persistence.RedisRepo{}
	err = repo.Connect("redis://" + s.Addr())

	if err != nil {
		panic(err)
	}

	h := GamesHandler{repo: &repo}
	req := httptest.NewRequest(http.MethodPost, "/v1/games", strings.NewReader(createGameReqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, h.createGame(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var game engine.Game

		json.Unmarshal(rec.Body.Bytes(), &game)

		assert.NotNil(t, game.ID)
		assert.Equal(t, "Rick", game.BoardSide1.Player.Name)
		assert.Equal(t, "Morty", game.BoardSide2.Player.Name)
	}
}

func TestUpdateGame(t *testing.T) {
	// Setup
	e := echo.New()
	g := engine.NewGame("Rick", "Morty")

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	repo := persistence.RedisRepo{}
	err = repo.Connect("redis://" + s.Addr())

	if err != nil {
		panic(err)
	}

	repo.SaveGame(g)

	h := GamesHandler{repo: &repo}
	req := httptest.NewRequest(http.MethodPatch, "/v1/games/:id", strings.NewReader(updateGameReqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/games/:id")
	c.SetParamNames("id")
	c.SetParamValues(g.ID.String())

	// Assertions
	if assert.NoError(t, h.updateGame(c)) {
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

func TestGetGame(t *testing.T) {
	// Setup
	e := echo.New()
	g := engine.NewGame("Rick", "Morty")

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	repo := persistence.RedisRepo{}
	err = repo.Connect("redis://" + s.Addr())

	if err != nil {
		panic(err)
	}

	repo.SaveGame(g)

	h := GamesHandler{repo: &repo}

	req := httptest.NewRequest(http.MethodGet, "/v1/games/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/v1/games/:id")
	c.SetParamNames("id")
	c.SetParamValues(g.ID.String())

	// Assertions
	if assert.NoError(t, h.getGame(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var game engine.Game

		json.Unmarshal(rec.Body.Bytes(), &game)

		assert.Equal(t, g.ID, game.ID)
		assert.Equal(t, "Rick", game.BoardSide1.Player.Name)
		assert.Equal(t, "Morty", game.BoardSide2.Player.Name)
	}
}
