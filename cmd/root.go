package cmd

import (
	"fmt"
	"os"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	. "github.com/arunvelsriram/sftp-exporter/pkg/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:   "sftp-exporter",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cfg)
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

	rootCmd.PersistentFlags().Int(FlagPort, 8080, "exporter port")
	viper.BindPFlag(ViperKeyPort, rootCmd.PersistentFlags().Lookup(FlagPort))

	rootCmd.PersistentFlags().String(FlagSFTPHost, "localhost", "sftp host")
	viper.BindPFlag(ViperKeySFTPHost, rootCmd.PersistentFlags().Lookup(FlagSFTPHost))

	rootCmd.PersistentFlags().Int(FlagSFTPPort, 22, "sftp port")
	viper.BindPFlag(ViperKeySFTPPort, rootCmd.PersistentFlags().Lookup(FlagSFTPPort))

	rootCmd.PersistentFlags().String(FlagSFTPUser, "", "sftp user")
	viper.BindPFlag(ViperKeySFTPUser, rootCmd.PersistentFlags().Lookup(FlagSFTPUser))

	rootCmd.PersistentFlags().String(FlagSFTPPass, "", "sftp user")
	viper.BindPFlag(ViperKeySFTPPass, rootCmd.PersistentFlags().Lookup(FlagSFTPPass))
}

func initConfig() {
	viper.AutomaticEnv()
	cfg = config.NewConfig()
}
