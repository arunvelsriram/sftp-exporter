package cmd

import (
	"fmt"
	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/arunvelsriram/sftp-exporter/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "sftp-exporter",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if err := server.Start(cfg); err != nil {
			log.WithFields(log.Fields{"event": "starting server"}).Fatal(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{"event": "executing command"}).Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String(c.FlagBindAddress, "127.0.0.1", "exporter bind address")
	_ = viper.BindPFlag(c.ViperKeyBindAddress, rootCmd.PersistentFlags().Lookup(c.FlagBindAddress))

	rootCmd.PersistentFlags().Int(c.FlagPort, 8080, "exporter port")
	_ = viper.BindPFlag(c.ViperKeyPort, rootCmd.PersistentFlags().Lookup(c.FlagPort))

	var levels = make([]string, len(log.AllLevels))
	for i, level := range log.AllLevels {
		levels[i] = level.String()
	}
	rootCmd.PersistentFlags().String(
		c.FlagLogLevel,
		log.InfoLevel.String(),
		fmt.Sprintf("log level [%s]", strings.Join(levels, " | ")),
	)
	_ = viper.BindPFlag(c.ViperKeyLogLevel, rootCmd.PersistentFlags().Lookup(c.FlagLogLevel))

	rootCmd.PersistentFlags().String(c.FlagSFTPHost, "localhost", "sftp host")
	_ = viper.BindPFlag(c.ViperKeySFTPHost, rootCmd.PersistentFlags().Lookup(c.FlagSFTPHost))

	rootCmd.PersistentFlags().Int(c.FlagSFTPPort, 22, "sftp port")
	_ = viper.BindPFlag(c.ViperKeySFTPPort, rootCmd.PersistentFlags().Lookup(c.FlagSFTPPort))

	rootCmd.PersistentFlags().String(c.FlagSFTPUser, "", "sftp user")
	_ = viper.BindPFlag(c.ViperKeySFTPUser, rootCmd.PersistentFlags().Lookup(c.FlagSFTPUser))

	rootCmd.PersistentFlags().String(c.FlagSFTPPass, "", "sftp user")
	_ = viper.BindPFlag(c.ViperKeySFTPPass, rootCmd.PersistentFlags().Lookup(c.FlagSFTPPass))
}

func initConfig() {
	viper.AutomaticEnv()
	cfg = config.LoadConfig()
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	level, err := log.ParseLevel(cfg.GetLogLevel())
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
}
