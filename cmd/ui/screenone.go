package ui

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(spotify.Playlist)
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

type screenOneModel struct {
	playlists list.Model
	selected  spotify.Playlist
	spinner   spinner.Model
	quitting  bool
}

func ScreenOne() *screenOneModel {
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

	// Spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &screenOneModel{
		playlists: l,
		spinner:   s,
	}
}

type resp string

func (s *screenOneModel) Init() tea.Cmd { return s.spinner.Tick }

func (s *screenOneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			fallthrough
		case "q":
			fallthrough
		case "esc":
			s.quitting = true
			return s, tea.Quit
		case "enter":
			if pl, ok := s.playlists.SelectedItem().(spotify.Playlist); ok {
				s.selected = pl
			}

			screenTwo := ScreenTwo(s.selected)

			// Bandaid fix, to identify the type sent from screen one
			return screenTwo.Update(resp(""))
		default:
			var cmd tea.Cmd
			s.playlists, cmd = s.playlists.Update(msg)
			return s, cmd
		}
	default:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}
}

func (s *screenOneModel) View() string {
	return "\n" + s.playlists.View()
}
