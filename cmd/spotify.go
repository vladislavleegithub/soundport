package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/Samarthbhat52/soundport/cmd/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	accent = lipgloss.NewStyle().Foreground(lipgloss.Color("163"))
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("84"))
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("161"))
)

func init() {
	rootCmd.AddCommand(spotifyCmd)
	spotifyCmd.AddCommand(spotifyLoginCmd)
	spotifyCmd.AddCommand(spotifyPlaylistsCmd)
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
	Use:    "login",
	Short:  "",
	Long:   "",
	Args:   cobra.NoArgs,
	PreRun: ensureInit,
	Run: func(cmd *cobra.Command, args []string) {
		var message strings.Builder
		var status strings.Builder

		creds := spotify.NewCredentials()

		message.WriteString("Click on " + accent.Render("Accept") + " in the browser popup\n")
		fmt.Println(message.String())

		ch := make(chan int)
		state := spotify.RandStringBytes(16)

		url := creds.GetAuthURL(state)
		go creds.StartHttpServer(ch, state)
		go spotify.OpenBrowser(url)

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

var spotifyPlaylistsCmd = &cobra.Command{
	Use:    "get",
	PreRun: ensureLogin,
	Run: func(cmd *cobra.Command, args []string) {
		a, _ := spotify.NewAuth()
		resp, err := a.GetPlaylists()
		if err != nil {
			log.Fatal(err)
		}

		playlistItems := resp.GetPlaylistItems()

		l := list.New(playlistItems, list.NewDefaultDelegate(), 0, 0)
		l.Title = "Select a playlist"
		m := ui.InitModel(l)

		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program: ", err)
			os.Exit(1)
		}
	},
}

// FIX : DECOUPLE SPOTIFY AND YT SETUP
func ensureInit(cmd *cobra.Command, args []string) {
	spfyId := viper.GetString("spfy-id")
	spfySecret := viper.GetString("spfy-secret")
	ytCookie := viper.GetString("yt-cookie")

	if spfyId == "" || spfySecret == "" || ytCookie == "" {
		fmt.Println("credentials not setup")
		fmt.Println("Please run `soundport setup` to initialize")
		os.Exit(1)
	}
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
