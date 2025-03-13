package spotify

import (
	"log"
	"math/rand"
	"net/url"

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
	base_auth_url  string = "https://accounts.spotify.com"
	auth_url       string = base_auth_url + "/authorize"
	access_tok_url string = base_auth_url + "/api/token"
	scope          string = "playlist-read-private playlist-read-collaborative"
	letterBytes    string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Auth struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
}

type Credentials struct {
	ClientId     string
	ClientSecret string
	State        string
}

func NewCredentials() *Credentials {
	clientId := viper.GetString("spfy-id")
	clientSecret := viper.GetString("spfy-secret")
	state := RandStringBytes(16)

	return &Credentials{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		State:        state,
	}
}

func (c *Credentials) GetAuthURL() string {
	u, err := url.Parse(auth_url)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.ClientId)
	q.Set("scope", scope)
	q.Set("redirect_uri", "http://"+redirect_url)
	q.Set("state", c.State)

	// Encode the Query
	u.RawQuery = q.Encode()

	return u.String()
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (c *Credentials) OpenBrowser(url string) {
	browser.Stdout = nil
	browser.Stderr = nil
	browser.OpenURL(url)
}
