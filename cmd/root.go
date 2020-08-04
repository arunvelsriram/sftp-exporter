package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/arunvelsriram/sftp-exporter/pkg/config"
	"github.com/mitchellh/go-homedir"
	"github.com/rclone/rclone/backend/sftp"
	"github.com/rclone/rclone/fs/config/obscure"
	"github.com/rclone/rclone/fs/operations"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "sftp-exporter",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		pass, err := obscure.Obscure("pass")
		if err != nil {
			panic(err)
		}
		sftpConfig := config.SFTPConfig{"host": "localhost", "port": "2220", "user": "basic", "pass": pass}
		fs, err := sftp.NewFs("my-sftp", "", sftpConfig)
		if err != nil {
			panic(err)
		}
		fmt.Println("got fs")

		size, count, err := operations.Count(context.Background(), fs)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Size: %d\n", size)
		fmt.Printf("Count: %d\n", count)
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
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configFileName := ".sftp-exporter.yaml"
	cfgFile = path.Join(home, configFileName)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file")
}

func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("error reading config file %s: %v", viper.ConfigFileUsed(), err))
	}
}
