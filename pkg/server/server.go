package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/gorilla/handlers"
)

func Start(cfg config.Config) error {
	r := http.NewServeMux()

	r.Handle("/healthz", withLogging(healthzHandler))

	addr := fmt.Sprintf("%s:%d", cfg.GetBindAddress(), cfg.GetPort())
	fmt.Printf("Listening on: %s\n", addr)
	return http.ListenAndServe(addr, r)
}

func withLogging(h http.HandlerFunc) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "healthy")
}
