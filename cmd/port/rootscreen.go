package port

import (
	"fmt"
	"os"

	"github.com/Samarthbhat52/soundport/api"
	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/Samarthbhat52/soundport/api/ytmusic"
	"github.com/Samarthbhat52/soundport/logger"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var glbLogger = logger.GetInstance()

type portModel struct {
	selectedPlId  string
	createdPlId   string
	statusMessage string
	src           api.Source
	dst           api.Destination
	playlists     list.Model // List of playlists in source provider
	quitting      bool
	successful    bool
}

type (
	playlistSelected struct{}
	playlistCreated  struct{}
)

func createSourceClient(srcFlag string) (api.Source, error) {
	switch srcFlag {
	case "spfy":
		return spotify.NewClient(), nil
	default:
		return nil, fmt.Errorf("invalid source flag '%s': only accepts 'spfy'", srcFlag)
	}
}

func createDestClient(destFlag string) (api.Destination, error) {
	switch destFlag {
	case "ytmusic":
		return ytmusic.NewClient(), nil
	default:
		return nil, fmt.Errorf("invalid destination flag '%s': only accepts 'ytmusic'", destFlag)
	}
}

func NewPortModel() *portModel {
	// Init source and destination
	var src api.Source
	var dest api.Destination

	src, err := createSourceClient(srcFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	dest, err = createDestClient(destFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
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

	return &portModel{
		src:       src,
		dst:       dest,
		playlists: l,
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

	if m.selectedPlId == "" {
		return m.updatePlaylists(msg)
	}

	return m.updatePortProgress(msg)
}

func (m portModel) View() string {
	var s string

	if m.selectedPlId == "" {
		s = m.viewPlaylists()
	} else {
		s = m.viewPortProgress()
	}

	return s
}
