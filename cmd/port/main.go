package port

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Samarthbhat52/soundport/api/spotify"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PortCmd = &cobra.Command{
	Use:    "port",
	PreRun: ensureLogin,
	Run: func(cmd *cobra.Command, args []string) {
		m := rootScreen()

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program: ", err)
			os.Exit(1)
		}
	},
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
