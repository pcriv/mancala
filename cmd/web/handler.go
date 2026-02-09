package main

import (
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/pcriv/mancala/proto"
	"github.com/pcriv/mancala/proto/protoconnect"
)

//go:embed templates/*.html
var templateFS embed.FS

type handler struct {
	client protoconnect.ServiceClient
	logger *slog.Logger
	tmpls  map[string]*template.Template
}

func newHandler(client protoconnect.ServiceClient, logger *slog.Logger) *handler {
	funcMap := template.FuncMap{
		"seq": func(n int) []int {
			s := make([]int, n)
			for i := range n {
				s[i] = i
			}

			return s
		},
		"reverse": func(s []int64) []int64 {
			r := make([]int64, len(s))
			for i, v := range s {
				r[len(s)-1-i] = v
			}

			return r
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"isPlayer1Turn": func(g *proto.Game) bool {
			return g.GetTurn() == proto.Turn_TURN_PLAYER1
		},
		"isPlayer2Turn": func(g *proto.Game) bool {
			return g.GetTurn() == proto.Turn_TURN_PLAYER2
		},
		"isGameOver": func(g *proto.Game) bool {
			return g.GetResult() != proto.Result_RESULT_UNSPECIFIED
		},
		"resultText": func(g *proto.Game) string {
			switch g.GetResult() {
			case proto.Result_RESULT_PLAYER1_WINS:
				return g.GetBoardSide1().GetPlayer().GetName() + " wins!"
			case proto.Result_RESULT_PLAYER2_WINS:
				return g.GetBoardSide2().GetPlayer().GetName() + " wins!"
			case proto.Result_RESULT_TIE:
				return "It's a tie!"
			default:
				return ""
			}
		},
	}

	shared := template.Must(
		template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/layout.html", "templates/board.html"),
	)

	tmpls := make(map[string]*template.Template)

	for _, page := range []string{"home.html", "game.html", "error.html"} {
		tmpls[page] = template.Must(template.Must(shared.Clone()).ParseFS(templateFS, "templates/"+page))
	}

	return &handler{
		client: client,
		logger: logger,
		tmpls:  tmpls,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", h.handleHome)
	mux.HandleFunc("POST /games", h.handleCreateGame)
	mux.HandleFunc("GET /games/{id}", h.handleGetGame)
	mux.HandleFunc("POST /games/{id}/play", h.handleExecutePlay)
	mux.HandleFunc("GET /games/{id}/board", h.handleGetBoard)
}

func (h *handler) renderError(w http.ResponseWriter, status int, title string, message string) {
	w.WriteHeader(status)

	data := map[string]string{
		"Title":   title,
		"Message": message,
	}

	if err := h.tmpls["error.html"].ExecuteTemplate(w, "layout", data); err != nil {
		h.logger.Error("failed to render error template", slog.String("error", err.Error()))
	}
}
