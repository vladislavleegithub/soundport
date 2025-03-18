package cmd

import (
	"fmt"
	"io"
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
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(spotify.Playlists)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

var portCmd = &cobra.Command{
	Use:    "port",
	PreRun: ensureLogin,
	Run: func(cmd *cobra.Command, args []string) {
		a, _ := spotify.NewAuth()

		playlists, err := a.GetPlaylists()
		if err != nil {
			log.Fatal(err)
		}
		plItems := playlists.GetPlaylistItems()

		l := list.New(plItems, itemDelegate{}, 20, 10)
		l.Title = "Choose a playlist to port from"
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = paginationStyle
		l.Styles.HelpStyle = helpStyle

		m := ui.InitModel(l)

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
