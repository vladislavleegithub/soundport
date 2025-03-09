package cmd

import (
	"fmt"

	"github.com/Samarthbhat52/soundport/api/ytmusic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(playlistCmd)
}

var playlistCmd = &cobra.Command{
	Use: "playlist",
	Run: func(cmd *cobra.Command, args []string) {
		err := ytmusic.PlaylistAdd([]string{"lKcSthXGNjY"})
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
