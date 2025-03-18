package spotify

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/Samarthbhat52/soundport/cmd/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
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

var getCmd = &cobra.Command{
	Use:    "get",
	PreRun: ensureLogin,
	Run: func(cmd *cobra.Command, args []string) {
		const defaultWidth = 20
		const lineHeight = 10

		a, _ := spotify.NewAuth()

		playlists, err := a.GetPlaylists()
		if err != nil {
			log.Fatal("Unable to fetch playlists: ", err)
		}

		plItems := playlists.GetPlaylistItems()

		l := list.New(plItems, itemDelegate{}, defaultWidth, lineHeight)
		m := ui.InitModel(l)

		_, err = tea.NewProgram(m).Run()
		if err != nil {
			fmt.Println("Error running program: ", err)
			os.Exit(1)
		}
	},
}
