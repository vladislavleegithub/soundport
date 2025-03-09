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
		vidIds, err := ytmusic.SearchSongYT(
			[]string{"eyes, bazzi", "Wrapped Around Your Finger, Post Malone"},
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = ytmusic.PlaylistAdd("Test", ytmusic.PUBLIC, vidIds)
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}
