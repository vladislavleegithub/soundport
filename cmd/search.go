package cmd

import (
	"fmt"
	"os"

	"github.com/Samarthbhat52/soundport/api/ytmusic"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use: "search",
	Run: func(cmd *cobra.Command, args []string) {
		visitorId, err := ytmusic.GetVisitorId()
		if err != nil {
			fmt.Println("ERROR FETCHING VISITOR ID: ", err)
			os.Exit(1)
		}

		fmt.Println("VISITOR ID: ", visitorId)
	},
}

/*
Body Format
{
  "context": {
    "client": {
      "hl": "en",
      "gl": "IN",
      "clientName": "WEB_REMIX",
      "clientVersion": "1.20250305.01.00",
    },
  },
  "query": <SONG_NAME>,
}
*/
