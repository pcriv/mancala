package main

import (
	"log/slog"
	"net/http"
	"strconv"

	"connectrpc.com/connect"

	"github.com/pcriv/mancala/proto"
)

func (h *handler) handleExecutePlay(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")

	pitIndexStr := r.FormValue("pit_index")

	pitIndex, parseErr := strconv.ParseInt(pitIndexStr, 10, 64)
	if parseErr != nil {
		h.renderBoardWithError(w, r, gameID, "Invalid pit index.")

		return
	}

	resp, err := h.client.ExecutePlay(r.Context(), connect.NewRequest(&proto.ExecutePlayRequest{
		GameId:   gameID,
		PitIndex: pitIndex,
	}))
	if err != nil {
		switch connect.CodeOf(err) {
		case connect.CodeNotFound:
			h.renderError(w, http.StatusNotFound, "Game Not Found", "The game you're looking for doesn't exist.")
		case connect.CodeInvalidArgument:
			h.renderBoardWithError(w, r, gameID, "Invalid move. Please try a different pit.")
		default:
			h.logger.ErrorContext(r.Context(), "failed to execute play", slog.String("error", err.Error()))
			h.renderError(w, http.StatusInternalServerError, "Error", "Failed to execute move. Please try again.")
		}

		return
	}

	data := map[string]any{
		"Game":    resp.Msg.GetGame(),
		"Animate": true,
	}

	if tmplErr := h.tmpls["game.html"].ExecuteTemplate(w, "board.html", data); tmplErr != nil {
		h.logger.ErrorContext(r.Context(), "failed to render board template", slog.String("error", tmplErr.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *handler) renderBoardWithError(w http.ResponseWriter, r *http.Request, gameID string, errorMsg string) {
	resp, err := h.client.FindGame(r.Context(), connect.NewRequest(&proto.FindGameRequest{
		Id: gameID,
	}))
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to find game for error render", slog.String("error", err.Error()))
		h.renderError(w, http.StatusInternalServerError, "Error", "Failed to load game.")

		return
	}

	data := map[string]any{
		"Game":  resp.Msg.GetGame(),
		"Error": errorMsg,
	}

	if tmplErr := h.tmpls["game.html"].ExecuteTemplate(w, "board.html", data); tmplErr != nil {
		h.logger.ErrorContext(r.Context(), "failed to render board template", slog.String("error", tmplErr.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
