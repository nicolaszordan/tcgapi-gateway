package router

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/nicolaszordan/tcgapi-gateway/internal/proxy"
)

type Config struct {
	LorcanaServiceURL string
	MTGServiceURL     string
}

type Router struct {
	config Config
	mux    *http.ServeMux
}

func NewRouter(cfg Config) *Router {
	mux := http.NewServeMux()

	r := &Router{
		config: cfg,
		mux:    mux,
	}

	mux.HandleFunc("/health", r.health)
	mux.HandleFunc("/cards", r.cardsHandler)

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	slog.Info("received request", "method", req.Method, "path", req.URL.Path)
	r.mux.ServeHTTP(w, req)
}

func (r *Router) health(w http.ResponseWriter, req *http.Request) {
	slog.Info("health check")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (r *Router) cardsHandler(w http.ResponseWriter, req *http.Request) {
	slog.Info("cards handler")

	// Forward request to the matching game service
	game := req.URL.Query().Get("game")
	resp, err := r.forwardRequestToService(game, req)
	if err != nil {
		http.Error(w, "failed to forward request to service", http.StatusBadRequest)
	}

	// Forward response from game service to client
	proxy.CopyResponse(w, resp)
}

func (r *Router) forwardRequestToService(game string, req *http.Request) (*http.Response, error) {
	switch game {
	case "lorcana":
		// Proxy to Lorcana service
		slog.Info("proxying to Lorcana service", "url", r.config.LorcanaServiceURL)
		return proxy.ForwardRequest(r.config.LorcanaServiceURL, req)
	case "mtg":
		// Proxy to MtG service
		slog.Info("proxying to MtG service", "url", r.config.MTGServiceURL)
		return proxy.ForwardRequest(r.config.MTGServiceURL, req)
	default:
		slog.Warn("unknown game parameter", "game", game)
		return nil, errors.New("unsupported game: " + game)
	}
}
