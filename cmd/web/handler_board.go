package main

import (
	"log/slog"
	"net/http"

	"connectrpc.com/connect"

	"github.com/pcriv/mancala/proto"
)

func (h *handler) handleGetBoard(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	data := map[string]any{
		"Game": resp.Msg.GetGame(),
	}

	if tmplErr := h.tmpls["game.html"].ExecuteTemplate(w, "board.html", data); tmplErr != nil {
		h.logger.ErrorContext(r.Context(), "failed to render board template", slog.String("error", tmplErr.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
