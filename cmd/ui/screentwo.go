package ui

import (
	"fmt"
	"os"

	"github.com/Samarthbhat52/soundport/api/spotify"
	"github.com/Samarthbhat52/soundport/api/ytmusic"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg          string
	responseMessage []string
	createPlaylist  bool
)

type playlistDetails struct {
	plName     string
	plDesc     string
	totalSongs int
}

type screenTwoModel struct {
	playlistDetails
	ch            chan []string
	songs         []string
	songIds       []string
	notFound      []string
	totalFound    int
	totalNotFound int
	quitting      bool
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

	plDetails := playlistDetails{
		totalSongs: tracks.Total,
		plName:     pl.Name,
		plDesc:     pl.Desc,
	}

	return &screenTwoModel{
		songs:           trackStrings,
		ch:              make(chan []string, tracks.Total),
		playlistDetails: plDetails,
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
			m.waitForActivity(),
			m.listenForActivity(m.songs),
		)

	case responseMessage:

		if msg[1] == "" {
			// TODO: Write not found songs onto a file
			m.totalNotFound++
			m.notFound = append(m.notFound, msg[0])
		} else {
			m.totalFound++
			m.songIds = append(m.songIds, msg[1])
		}

		if m.totalNotFound+m.totalFound == m.totalSongs {
			return m, m.createPlaylistFunc()
		}

		return m, m.waitForActivity()

	case createPlaylist:
		m.quitting = true
		return m, tea.Quit

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

func (m *screenTwoModel) waitForActivity() tea.Cmd {
	return func() tea.Msg {
		return responseMessage(<-m.ch)
	}
}

func (m *screenTwoModel) listenForActivity(songs []string) tea.Cmd {
	return func() tea.Msg {
		ytmusic.GetSongInfo(songs, m.ch)
		return nil
	}
}

func (m *screenTwoModel) createPlaylistFunc() tea.Cmd {
	return func() tea.Msg {
		err := ytmusic.PlaylistAdd(m.plName, "PRIVATE", m.songIds)
		if err != nil {
			glbLogger.Println("error creating playlist: ", err)
			os.Exit(1)
		}

		return createPlaylist(true)
	}
}

func constructSongsList(track spotify.PlaylistTracks) []string {
	songsList := make([]string, track.Total)

	for i := range track.Total {
		sn := track.Tracks[i].Track.Name
		sn += " By: " + track.Tracks[i].Track.Artists[0].Name
		songsList[i] = sn
	}

	return songsList
}
