package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/arunvelsriram/sftp-exporter/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		fmt.Println(err)
		os.Exit(1)
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

	rootCmd.PersistentFlags().String(c.FlagSFTPPass, "", "sftp password")
	_ = viper.BindPFlag(c.ViperKeySFTPPass, rootCmd.PersistentFlags().Lookup(c.FlagSFTPPass))

	rootCmd.PersistentFlags().String(c.FlagSFTPKey, "", "sftp key (base64 encoded)")
	_ = viper.BindPFlag(c.ViperKeySFTPKey, rootCmd.PersistentFlags().Lookup(c.FlagSFTPKey))

	rootCmd.PersistentFlags().String(c.FlagSFTPKeyFile, "", "sftp key file")
	_ = viper.BindPFlag(c.ViperKeySFTPKeyFile, rootCmd.PersistentFlags().Lookup(c.FlagSFTPKeyFile))

	rootCmd.PersistentFlags().String(c.FlagSFTPKeyPassphrase, "", "sftp key passphrase")
	_ = viper.BindPFlag(c.ViperKeySFTPKeyPassphrase, rootCmd.PersistentFlags().Lookup(c.FlagSFTPKeyPassphrase))
}

func initConfig() {
	viper.AutomaticEnv()

	var err error
	fs := afero.NewOsFs()
	cfg = config.MustLoadConfig(fs)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	level, err := log.ParseLevel(cfg.GetLogLevel())
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)
}
