package cmd

import (
	"fmt"
	"os"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	c "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/arunvelsriram/sftp-exporter/pkg/server"
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
			fmt.Println(err)
			os.Exit(1)
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
}
