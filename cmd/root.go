package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const CONFIG_FILE_NAME = ".soundport.json"

var (
	cfgFile = ""
	rootCmd = &cobra.Command{
		Use:   "soundport",
		Short: "",
		Long:  ``,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
}

func initConfig() {
	viper.SetDefault("yt-Cookie", "")
	viper.SetDefault("yt-auth-token", "")
	viper.SetDefault("spfy-id", "")
	viper.SetDefault("spfy-secret", "")
	viper.SetDefault("spfy-access", "")
	viper.SetDefault("spfy-refresh", "")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		err := viper.ReadInConfig()
		cobra.CheckErr(err)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// default file path.
		defaultPath := path.Join(home, CONFIG_FILE_NAME)

		// set viper config path
		viper.SetConfigFile(defaultPath)
		err = viper.ReadInConfig()
		if err != nil {
			// Create a file
			viper.SafeWriteConfigAs(defaultPath)
			err = viper.ReadInConfig()
			cobra.CheckErr(err)
		}
	}
}
