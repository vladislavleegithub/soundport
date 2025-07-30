package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vladislavleegithub/soundport/cmd/port"
	"github.com/vladislavleegithub/soundport/cmd/spotify"
	"github.com/vladislavleegithub/soundport/cmd/ytmusic"
)

const CONFIG_FILE_NAME = ".soundport.json"

var rootCmd = &cobra.Command{
	Use:   "soundport",
	Short: "The root command for Soundport CLI, used to manage and interact with music services.",
}

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
