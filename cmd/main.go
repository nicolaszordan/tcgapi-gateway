package main

import (
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("tcgapi-gateway starting")

	http.HandleFunc("/health", health)

	slog.Info("starting server on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		slog.Error("server error", "error", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	slog.Info("health check")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
