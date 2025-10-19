package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/nicolaszordan/tcgapi-gateway/internal/router"
)

type Config struct {
	BindPort     string
	RouterConfig router.Config
}

func main() {
	slog.Info("tcgapi-gateway starting")

	config := loadConfig()

	r := router.NewRouter(config.RouterConfig)

	slog.Info("starting server on :" + config.BindPort)
	err := http.ListenAndServe(":"+config.BindPort, r)
	if err != nil {
		slog.Error("server error", "error", err)
	}
}

func loadConfig() Config {
	return Config{
		BindPort: "8000",
		RouterConfig: router.Config{
			LorcanaServiceURL: os.Getenv("LORCANA_SERVICE_URL"),
			MTGServiceURL:     os.Getenv("MTG_SERVICE_URL"),
		},
	}
}
