package port

import (
	"fmt"
	"os"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/Samarthbhat52/soundport/api/types"
	"github.com/Samarthbhat52/soundport/api/ytmusic"
	"github.com/Samarthbhat52/soundport/logger"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

var glbLogger = logger.GetInstance()

type portModel struct {
	songs     []string
	ch        chan types.SongDetails
	selected  *types.PlaylistDetails
	src       types.Source
	dst       types.Destination
	playlists list.Model
	progress  progress.Model
	total     int
	percent   float64
	quitting  bool
}

type (
	playlistSelected   struct{}
	playlistCreated    struct{}
	songSearchComplete bool
)

func NewPortModel() *portModel {
	// Init source and destination
	var src types.Source
	var dest types.Destination

	switch srcFlag {
	default:
		src = spotify.NewClient()
	}

	switch destFlag {
	default:
		dest = ytmusic.NewClient()
	}

	// Get source playlists
	playlists, err := src.GetPlaylists()
	if err != nil {
		glbLogger.Printf("Error getting playlists: %s", err)
		fmt.Println("Something went wrong")
		os.Exit(1)
	}

	// Init list
	l := list.New(playlists, list.NewDefaultDelegate(), 20, 10)
	l.Title = "Choose a playlist"

	// Init progress bar
	prog := progress.New(progress.WithDefaultGradient())

	return &portModel{
		src:       src,
		dst:       dest,
		playlists: l,
		progress:  prog,
	}
}

func (m portModel) Init() tea.Cmd {
	return nil
}

func (m portModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	if m.selected == nil {
		return m.updatePlaylists(msg)
	}

	return m.updatePortProgress(msg)
}

func (m portModel) View() string {
	var s string

	if m.quitting {
		return "Something to show when quitting"
	}

	if m.selected == nil {
		s = m.viewPlaylists()
	} else {
		s = m.viewPortProgress()
	}

	return s
}
