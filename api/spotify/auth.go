package spotify

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Samarthbhat52/soundport/logger"
	"github.com/spf13/viper"
)

var glbLogger = logger.GetInstance()

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
	AuthUrl      string
}

// Used for making authenticated queries to spotify api
type auth struct {
	accessToken string
	playlistUrl string
}

func NewCredentials() *credentials {
	clientId := viper.GetString("spfy-id")
	clientSecret := viper.GetString("spfy-secret")
	state := randStringBytes(16)
	url := getAuthURL(clientId, state)

	return &credentials{
		clientId:     clientId,
		clientSecret: clientSecret,
		state:        state,
		AuthUrl:      url,
	}
}

func NewClient() *auth {
	accessToken := viper.GetString("spfy-access")

	return &auth{
		accessToken: accessToken,
	}
}

func getAccessToken(authCode string) error {
	body := url.Values{}
	body.Add("code", authCode)
	body.Add("redirect_uri", "http://"+redirect_url)
	body.Add("grant_type", "authorization_code")
	encodedBody := strings.NewReader(body.Encode())

	err := makeAuthRequest(encodedBody)
	if err != nil {
		return err
	}

	return nil
}

func RefreshSession() error {
	body := url.Values{}
	body.Add("grant_type", "refresh_token")
	body.Add("refresh_token", viper.GetString("spfy-refresh"))
	encodedBody := strings.NewReader(body.Encode())

	err := makeAuthRequest(encodedBody)
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

func (c *credentials) StartHttpServer(ch chan int) error {
	handleCallback := func(w http.ResponseWriter, r *http.Request) {
		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			glbLogger.Println("Error decoding query params: ", err.Error())
			ch <- -1
		}

		error := params.Get("error")
		if error != "" {
			glbLogger.Println("Permission denied: ", error)
			ch <- -1
		}

		// Check state sent from spotify
		retState := params.Get("state")
		if retState != c.state {
			glbLogger.Println("State mismatch error")
			ch <- -1
		}

		authCode := params.Get("code")
		if authCode == "" {
			glbLogger.Println("No auth token: ", authCode)
			ch <- -1
		}

		// Get access_token and refresh_token
		err = getAccessToken(authCode)
		if err != nil {
			glbLogger.Println("Error getting access token: ", err)
			ch <- -1
		}
		ch <- 0
	}

	http.HandleFunc("/callback", handleCallback)

	return http.ListenAndServe(server_url, nil)
}
