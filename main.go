package main

import "github.com/arunvelsriram/sftp-exporter/cmd"

//go:generate mkdir -p pkg/internal/mocks
//go:generate mockgen -source pkg/client/sftp_client.go -destination pkg/internal/mocks/sftp_client.go -package mocks

var version string = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
