package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus"
)

func Start(cfg config.Config) error {
	r := http.NewServeMux()

	r.Handle("/healthz", healthzHandler())
	r.Handle("/metrics", metricsHandler())

	addr := fmt.Sprintf("%s:%d", cfg.GetBindAddress(), cfg.GetPort())
	fmt.Printf("Listening on: %s\n", addr)
	return http.ListenAndServe(addr, r)
}

func withLogging(h http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, h)
}

func healthzHandler() http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "healthy")
	})
	return withLogging(handler)
}

func metricsHandler() http.Handler {
	handler := prometheus.Handler()
	return withLogging(handler)
}
