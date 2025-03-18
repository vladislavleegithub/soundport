package spotify

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	Cmd.AddCommand(setupCmd)
	Cmd.AddCommand(loginCmd)
	Cmd.AddCommand(getCmd)
}

var Cmd = &cobra.Command{
	Use:   "spotify",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
}

func ensureLogin(cmd *cobra.Command, args []string) {
	spfyAccess := viper.GetString("spfy-access")
	spfyRefresh := viper.GetString("spfy-refresh")

	if spfyAccess == "" || spfyRefresh == "" {
		fmt.Println("Not logged into spotify")
		fmt.Println("Please run `soundport spotify login`")
		os.Exit(1)
	}

	expiresAt := viper.GetTime("spfy-expires-at")

	// Check if auth token is close to expiry
	checkTime := expiresAt.Add(-10 * time.Minute)
	if time.Now().Before(checkTime) {
		return
	}

	err := spotify.RefreshSession()
	if err != nil {
		log.Fatal("error refreshing session")
	}
}
