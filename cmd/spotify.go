package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Samarthbhat52/soundport/api/spotify"
	textinputs "github.com/Samarthbhat52/soundport/cmd/ui/textInputs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var green = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))

func init() {
	rootCmd.AddCommand(spotifyCmd)
	spotifyCmd.AddCommand(spotifyInitCmd)
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
}

var spotifyInitCmd = &cobra.Command{
	Use:   "init",
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
		var message strings.Builder
		var status strings.Builder

		creds := spotify.NewCredentials()

		message.WriteString("Click on " + green.Render("Accept") + " in the browser popup\n")
		fmt.Println(message.String())

		ch := make(chan int)
		url := creds.GetAuthURL()
		go creds.StartHttpServer(ch)
		go creds.OpenBrowser(url)

		val := <-ch
		if val == 0 {
			status.WriteString(green.Render("Login successful\n"))
			fmt.Println(status.String())
		}
		fmt.Println("Browser window/tab can be closed.")
	},
}
