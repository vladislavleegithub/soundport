package spotify

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/Samarthbhat52/soundport/cmd/ui"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:    "login",
	Short:  "Authenticates the user with Spotify.",
	Long:   "This command logs the user into their Spotify account. If credentials are already stored, it refreshes the authentication tokens to maintain access to Spotify's services.",
	Args:   cobra.NoArgs,
	PreRun: ensureInit,
	Run: func(cmd *cobra.Command, args []string) {
		var message strings.Builder
		var status strings.Builder

		creds := spotify.NewCredentials()

		message.WriteString("Click on " + ui.Accent.Render("Accept") + " in the browser popup\n")
		fmt.Println(message.String())

		ch := make(chan int)
		// Handles callback
		go creds.StartHttpServer(ch)
		go creds.OpenBrowser()
		val := <-ch

		if val == 0 {
			status.WriteString(ui.Green.Render("Login successful\n"))
			fmt.Println(status.String())
		} else {
			status.WriteString(ui.Red.Render("Login failed\n"))
			fmt.Println(status.String())
		}
		fmt.Printf("Browser window/tab can be closed.\n\n")
	},
}

func ensureInit(cmd *cobra.Command, args []string) {
	spfyId := viper.GetString("spfy-id")
	spfySecret := viper.GetString("spfy-secret")

	if spfyId == "" || spfySecret == "" {
		fmt.Println("spotify credentials not setup")
		fmt.Println("Please run `soundport spotify setup`")
		os.Exit(1)
	}
}
