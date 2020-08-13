package main

import "github.com/arunvelsriram/sftp-exporter/cmd"

//go:generate mkdir -p pkg/internal/mocks
//go:generate mockgen -source pkg/config/config.go -destination pkg/internal/mocks/config.go -package mocks
func main() {
	cmd.Execute()
}
