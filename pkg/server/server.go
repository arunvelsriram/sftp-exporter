package server

import (
	"fmt"
	"net/http"

	"github.com/arunvelsriram/sftp-exporter/pkg/service"

	"github.com/arunvelsriram/sftp-exporter/pkg/client"
	"github.com/arunvelsriram/sftp-exporter/pkg/collector"
	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func Start(cfg config.Config) error {
	sftpClient := client.NewSFTPClient(cfg)
	sftpService := service.NewSFTPService(cfg, sftpClient)
	sftpCollector := collector.NewSFTPCollector(cfg, sftpService)
	prometheus.MustRegister(sftpCollector)

	r := http.NewServeMux()
	r.Handle("/healthz", WithLogging(healthzHandler()))
	r.Handle("/metrics", WithLogging(promhttp.Handler()))

	addr := fmt.Sprintf("%s:%d", cfg.GetBindAddress(), cfg.GetPort())
	log.Infof("will be listening on: %s", addr)
	return http.ListenAndServe(addr, r)
}

func healthzHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "healthy")
	}
	return http.HandlerFunc(fn)
}
