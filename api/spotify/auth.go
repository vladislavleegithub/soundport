package spotify

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

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

type authResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
	Scope        string        `json:"scope"`
	ExpiresIn    time.Duration `json:"expires_in"`
}

type credentials struct {
	clientId     string
	clientSecret string
}

// Used for making authenticated queries to spotify api
type auth struct {
	accessToken  string
	refreshToken string
	expiresAt    time.Time
	creds        *credentials
}

func NewCredentials() *credentials {
	clientId := viper.GetString("spfy-id")
	clientSecret := viper.GetString("spfy-secret")

	return &credentials{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func NewAuth() (*auth, error) {
	accessToken := viper.GetString("spfy-access")
	refreshToken := viper.GetString("spfy-refresh")
	expiresAt := viper.GetTime("spfy-expires-at")

	creds := NewCredentials()
	if accessToken == "" || refreshToken == "" || expiresAt.IsZero() {
		// TODO: Improve error message
		return nil, errors.New("not logged in")
	}

	return &auth{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
		creds:        creds,
	}, nil
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func OpenBrowser(url string) {
	browser.Stdout = nil
	browser.Stderr = nil
	browser.OpenURL(url)
}

func updateSpotifyConfig(resp *authResponse) error {
	// Set access and refresh_token to viper config
	viper.Set("spfy-access", resp.AccessToken)
	viper.Set("spfy-refresh", resp.RefreshToken)

	// Add token expiry
	expiresAt := time.Now().Add(resp.ExpiresIn * time.Second)
	viper.Set("spfy-expires-at", expiresAt)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func handleAuthResponse(req *http.Request) error {
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("error processing auth request")
	}

	authResponse := authResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&authResponse)
	if err != nil {
		return err
	}

	err = updateSpotifyConfig(&authResponse)
	if err != nil {
		return err
	}

	return nil
}
