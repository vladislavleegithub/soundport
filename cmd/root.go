package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/Samarthbhat52/soundport/cmd/port"
	"github.com/Samarthbhat52/soundport/cmd/spotify"
	"github.com/Samarthbhat52/soundport/cmd/ytmusic"
	"github.com/Samarthbhat52/soundport/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const CONFIG_FILE_NAME = ".soundport.json"

var (
	glbLogger = logger.GetInstance()
	rootCmd   = &cobra.Command{
		Use: "soundport",
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
	rootCmd.AddCommand(port.PortCmd)
	rootCmd.AddCommand(spotify.Cmd)
	rootCmd.AddCommand(ytmusic.Cmd)
}

func initConfig() {
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
