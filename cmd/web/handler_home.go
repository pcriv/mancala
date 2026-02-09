package main

import (
	"log/slog"
	"net/http"
)

func (h *handler) handleHome(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpls["home.html"].ExecuteTemplate(w, "layout", map[string]string{}); err != nil {
		h.logger.ErrorContext(r.Context(), "failed to render home template", slog.String("error", err.Error()))
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
