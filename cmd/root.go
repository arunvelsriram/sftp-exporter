package cmd

import (
	"fmt"
	"strings"

	"github.com/arunvelsriram/sftp-exporter/pkg/constants/viperkeys"
	"github.com/arunvelsriram/sftp-exporter/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "sftp-exporter",
	Short: "Prometheus Exporter for SFTP.",
	Run: func(cmd *cobra.Command, args []string) {
		level, err := log.ParseLevel(viper.GetString(viperkeys.LogLevel))
		if err != nil {
			log.Fatalf("Failed to set log level: %v", err)
		}
		log.SetLevel(level)

		log.Debugf("All configs:")
		for key, value := range viper.AllSettings() {
			if key == viperkeys.SFTPPassword || key == viperkeys.SFTPKey || key == viperkeys.SFTPKeyPassphrase {
				value = "**********"
			}
			log.Debugf("%s: %v", key, value)
		}

		if err = server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v\n", err)
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&configFile, viperkeys.ConfigFile, "c", "sftp-exporter.yaml", "exporter config file")
	rootCmd.Flags().String(viperkeys.BindAddress, "127.0.0.1", "exporter bind address")
	rootCmd.Flags().Int(viperkeys.Port, 8080, "exporter port")
	var logLevels = make([]string, len(log.AllLevels))
	for i, level := range log.AllLevels {
		logLevels[i] = level.String()
	}
	logLevelUsage := fmt.Sprintf("log level [%s]", strings.Join(logLevels, " | "))

	rootCmd.Flags().String(viperkeys.LogLevel, log.InfoLevel.String(), logLevelUsage)
	rootCmd.Flags().String(viperkeys.SFTPHost, "localhost", "SFTP host")
	rootCmd.Flags().Int(viperkeys.SFTPPort, 22, "SFTP port")
	rootCmd.Flags().String(viperkeys.SFTPUser, "", "SFTP user")
	rootCmd.Flags().String(viperkeys.SFTPPassword, "", "SFTP password")
	rootCmd.Flags().String(viperkeys.SFTPKey, "", "SFTP key (base64 encoded)")
	rootCmd.Flags().String(viperkeys.SFTPKeyPassphrase, "", "SFTP key passphrase")
	rootCmd.Flags().StringSlice(viperkeys.SFTPPaths, []string{"/"}, "SFTP paths")

	err := viper.BindPFlags(rootCmd.Flags())
	if err != nil {
		log.Fatalf("Viper failed to bind flags: %v", err)
	}
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // replace - in cmdline flags to _ in env vars
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Warnf("Failed to load config file %s: %v", configFile, err)
	}
	viper.AutomaticEnv()
}
