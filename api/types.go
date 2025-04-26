package api

type PlaylistDetails struct {
	PlId        string
	PlName      string
	PlDesc      string
	TotalTracks int
}

func (p PlaylistDetails) FilterValue() string { return p.PlName }

type SongDetails struct {
	Name  string
	Id    string
	Found bool
}
