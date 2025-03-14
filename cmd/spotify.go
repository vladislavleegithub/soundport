package cmd

import (
	"fmt"
	"strings"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	accent = lipgloss.NewStyle().Foreground(lipgloss.Color("163"))
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("161"))
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

		message.WriteString("Click on " + accent.Render("Accept") + " in the browser popup\n")
		fmt.Println(message.String())

		ch := make(chan int)
		url := creds.GetAuthURL()
		go creds.StartHttpServer(ch)
		go creds.OpenBrowser(url)

		val := <-ch
		if val == 0 {
			status.WriteString(green.Render("Login successful\n"))
			fmt.Println(status.String())
		} else {
			status.WriteString(red.Render("Login failed\n"))
			fmt.Println(status.String())
		}
		fmt.Println("Browser window/tab can be closed.")
	},
}
