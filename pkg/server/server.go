package server

import (
	"fmt"
	"net/http"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/gin-gonic/gin"
)

func Start(cfg config.Config) {
	router := gin.Default()

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "healthy")
	})

	router.Run(fmt.Sprintf("%s:%d", cfg.GetBindAddress(), cfg.GetPort()))
}
