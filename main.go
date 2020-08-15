package main

import "github.com/arunvelsriram/sftp-exporter/cmd"

//go:generate mkdir -p pkg/internal/mocks
//go:generate mockgen -source pkg/config/config.go -destination pkg/internal/mocks/config.go -package mocks
//go:generate mockgen -source pkg/client/sftp_client.go -destination pkg/internal/mocks/sftp_client.go -package mocks
//go:generate mockgen -source pkg/service/sftp_service.go -destination pkg/internal/mocks/sftp_service.go -package mocks
func main() {
	cmd.Execute()
}
