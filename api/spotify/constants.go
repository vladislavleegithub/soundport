package spotify

const (
	base_url     string = "127.0.0.1"
	port         string = "4214"
	server_url   string = base_url + ":" + port
	redirect_url string = server_url + "/callback"
)

const (
	// Auth urls
	base_auth_url string = "https://accounts.spotify.com"
	auth_url      string = base_auth_url + "/authorize"
	token_url     string = base_auth_url + "/api/token"
	scope         string = "playlist-read-private playlist-read-collaborative"

	// Other urls
	base_api_url string = "https://api.spotify.com/v1"
	playlist_url string = base_api_url + "/me/playlists"

	// used to generate secret key
	letterBytes string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)
