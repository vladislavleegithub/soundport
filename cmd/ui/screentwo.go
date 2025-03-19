package ui

import (
	"fmt"
	"os"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/Samarthbhat52/soundport/api/ytmusic"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg string

type screenTwoModel struct {
	songs      []string
	totalSongs int
	totalFound int
	ch         chan string
	quitting   bool
	err        errMsg
}

func ScreenTwo(pl spotify.Playlist) *screenTwoModel {
	a, _ := spotify.NewAuth()
	songs, err := a.GetTracks(pl.Tracks.Link)
	if err != nil {
		glbLogger.Printf("error getting playlist tracks: %s", err)
		fmt.Println("error getting playlist tracks")
		os.Exit(1)
	}
	tracks := spotify.PlaylistTracks(songs)

	trackStrings := constructSongsList(tracks)

	return &screenTwoModel{
		songs:      trackStrings,
		totalSongs: tracks.Total,
		ch:         make(chan string, tracks.Total),
	}
}

func (m *screenTwoModel) Init() tea.Cmd {
	return nil
}

func (m *screenTwoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			fallthrough
		case "q":
			fallthrough
		case "esc":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case resp:
		return m, tea.Batch(
			waitForActivity(m.ch),
			listenForActivity(m.ch, m.songs),
		)

	case responseMessage:
		m.totalFound++
		if m.totalSongs == m.totalFound {
			m.quitting = true
			return m, tea.Quit
		}
		return m, waitForActivity(m.ch)

	default:
		return m, nil
	}
}

func (m *screenTwoModel) View() string {
	s := fmt.Sprintf(
		"\n Songs found: %d\n\nPress `q` to exit\n",
		m.totalFound,
	)
	if m.quitting {
		s += "\n"
	}

	return s
}

type responseMessage string

func waitForActivity(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return responseMessage(<-sub)
	}
}

func listenForActivity(sub chan string, songs []string) tea.Cmd {
	return func() tea.Msg {
		for _, sn := range songs {
			go ytmusic.GetSongInfo(sn, sub)
		}

		return nil
	}
}

func constructSongsList(track spotify.PlaylistTracks) []string {
	songsList := make([]string, track.Total)

	for i := range track.Total {
		sn := track.Tracks[i].Track.Name
		sn += " " + track.Tracks[i].Track.Artists[0].Name
		sn += " " + track.Tracks[i].Track.Album.Name
		songsList[i] = sn
	}

	return songsList
}
