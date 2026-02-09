package main

import (
	"log/slog"
	"net/http"

	"connectrpc.com/connect"

	"github.com/pcriv/mancala/proto"
)

func (h *handler) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	player1 := r.FormValue("player1")
	player2 := r.FormValue("player2")

	if player1 == "" || player2 == "" {
		data := map[string]string{
			"Error":   "Both player names are required",
			"Player1": player1,
			"Player2": player2,
		}

		if tmplErr := h.tmpls["home.html"].ExecuteTemplate(w, "layout", data); tmplErr != nil {
			h.logger.ErrorContext(r.Context(), "failed to render home template", slog.String("error", tmplErr.Error()))
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}

		return
	}

	resp, err := h.client.CreateGame(r.Context(), connect.NewRequest(&proto.CreateGameRequest{
		Player1: player1,
		Player2: player2,
	}))
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to create game", slog.String("error", err.Error()))
		h.renderError(w, http.StatusInternalServerError, "Error", "Failed to create game. Please try again.")

		return
	}

	http.Redirect(w, r, "/games/"+resp.Msg.GetCreatedGame().GetId(), http.StatusSeeOther)
}

func (h *handler) handleGetGame(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")

	resp, err := h.client.FindGame(r.Context(), connect.NewRequest(&proto.FindGameRequest{
		Id: gameID,
	}))
	if err != nil {
		if connect.CodeOf(err) == connect.CodeNotFound {
			h.renderError(w, http.StatusNotFound, "Game Not Found", "The game you're looking for doesn't exist.")

			return
		}

		h.logger.ErrorContext(r.Context(), "failed to find game", slog.String("error", err.Error()))
		h.renderError(w, http.StatusInternalServerError, "Error", "Failed to load game. Please try again.")

		return
	}

	data := map[string]any{
		"Game": resp.Msg.GetGame(),
	}

	if tmplErr := h.tmpls["game.html"].ExecuteTemplate(w, "layout", data); tmplErr != nil {
		h.logger.ErrorContext(r.Context(), "failed to render game template", slog.String("error", tmplErr.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
