package spotify

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
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
	state        string
}

// Used for making authenticated queries to spotify api
type auth struct {
	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

func NewCredentials() *credentials {
	clientId := viper.GetString("spfy-id")
	clientSecret := viper.GetString("spfy-secret")
	state := RandStringBytes(16)

	return &credentials{
		clientId:     clientId,
		clientSecret: clientSecret,
		state:        state,
	}
}

func NewAuth() (*auth, error) {
	accessToken := viper.GetString("spfy-access")
	refreshToken := viper.GetString("spfy-refresh")
	expiresAt := viper.GetTime("spfy-expires-at")

	return &auth{
		accessToken:  accessToken,
		refreshToken: refreshToken,
		expiresAt:    expiresAt,
	}, nil
}

func getAccessToken(authCode string) (string, error) {
	body := url.Values{}
	body.Add("code", authCode)
	body.Add("redirect_uri", "http://"+redirect_url)
	body.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", token_url, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}

	err = handleAuthResponse(req)
	if err != nil {
		return "", err
	}

	return "processed", nil
}

func RefreshSession() error {
	fmt.Println("REFRESH SESSION CALLED")

	// Set up request body
	body := url.Values{}
	body.Add("grant_type", "refresh_token")
	body.Add("refresh_token", viper.GetString("spfy-refresh"))
	encodedBody := strings.NewReader(body.Encode())

	req, err := http.NewRequest("POST", token_url, encodedBody)
	if err != nil {
		return err
	}

	err = handleAuthResponse(req)
	if err != nil {
		return err
	}

	return nil
}

func getAuthorizationHeader() string {
	clientId := viper.GetString("spfy-id")
	clientSecret := viper.GetString("spfy-secret")

	data := fmt.Appendf([]byte{}, "%s:%s", clientId, clientSecret)

	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))

	base64.StdEncoding.Encode(dst, data)

	return fmt.Sprintf("Basic %s", string(dst))
}

func (c *credentials) GetAuthURL() string {
	u, err := url.Parse(auth_url)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.clientId)
	q.Set("scope", scope)
	q.Set("redirect_uri", "http://"+redirect_url)
	q.Set("state", c.state)

	// Encode the Query
	u.RawQuery = q.Encode()

	return u.String()
}
