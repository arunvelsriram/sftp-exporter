package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arunvelsriram/sftp-exporter/pkg/collector"
	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Start(cfg config.Config) error {
	sftpCollector := collector.NewSFTPCollector(cfg)
	prometheus.MustRegister(sftpCollector)

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
	handler := promhttp.Handler()
	return withLogging(handler)
}
