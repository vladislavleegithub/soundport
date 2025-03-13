package cmd

import (
	"fmt"
	"log"

	"github.com/Samarthbhat52/soundport/api/spotify"
	textinputs "github.com/Samarthbhat52/soundport/cmd/ui/textInputs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(spotifyCmd)
	spotifyCmd.AddCommand(spotifyLoginCmd)
}

type listOptions struct {
	options []string
}

var spotifyCmd = &cobra.Command{
	Use:   "spotify",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(textinputs.InitialModel())
		_, err := p.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var spotifyLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// post it onto screen.
		// Create a channel, to accept the status code.
		// Start a server.
		// if status is err, immediately add it to the channel and shut down server
		// if status is ok, add access and refresh tokens to config file.
		// shut down server

		creds := spotify.NewCredentials()
		authUrl := creds.GetAuthURL()
		fmt.Printf(
			"Please click the link below to sign in:\n\n%s\n\nclick accept and close the browser once done.\n",
			authUrl,
		)

		ch := make(chan string)
		go creds.StartHttpServer(ch)

		<-ch
	},
}
