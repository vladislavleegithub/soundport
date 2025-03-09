package cmd

import (
	"github.com/Samarthbhat52/soundport/api/ytmusic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use: "search",
	Run: func(cmd *cobra.Command, args []string) {
		ytmusic.SearchSongYT([]string{"eyes, bazzi", "Wrapped Around Your Finger, Post Malone"})
	},
}
