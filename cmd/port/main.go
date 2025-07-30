package port

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vladislavleegithub/soundport/api/spotify"
)

var (
	srcFlag  string
	destFlag string
)

var PortCmd = &cobra.Command{
	Use:    "port",
	Short:  "Select and port a playlist from Spotify to YT Music.",
	Long:   "Promts user to select a playlist from the ones found in their account. When selected, it starts the process of porting all the tracks found in that playlist.",
	PreRun: ensureLogin,
	Run: func(cmd *cobra.Command, args []string) {
		m := NewPortModel()

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	PortCmd.Flags().
		StringVarP(&srcFlag, "src", "s", "", "Source platform to move the playlist from")
	PortCmd.Flags().
		StringVarP(&destFlag, "dest", "d", "", "Destination platform to move the playlist to")
}

func ensureLogin(cmd *cobra.Command, args []string) {
	spfyAccess := viper.GetString("spfy-access")
	spfyRefresh := viper.GetString("spfy-refresh")

	if spfyAccess == "" || spfyRefresh == "" {
		fmt.Println("Not logged into spotify")
		fmt.Println("Please run `soundport spotify login`")
		os.Exit(1)
	}

	ytCookie := viper.GetString("yt-cookie")
	if ytCookie == "" {
		fmt.Println("Not setup youtube")
		fmt.Println("Please run `soundport ytmusic setup`")
		os.Exit(1)
	}

	expiresAt := viper.GetTime("spfy-expires-at")

	// Refresh even if auth token is close to expiry
	checkTime := expiresAt.Add(-10 * time.Minute)
	if time.Now().Before(checkTime) {
		return
	}

	err := spotify.RefreshSession()
	if err != nil {
		log.Fatal("error refreshing session")
	}
}
