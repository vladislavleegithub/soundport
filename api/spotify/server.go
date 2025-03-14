package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Samarthbhat52/soundport/logger"
	"github.com/spf13/viper"
)

var glbLogger = logger.GetInstance()

func (c *Credentials) StartHttpServer(ch chan int) error {
	handleCallback := func(w http.ResponseWriter, r *http.Request) {
		// Unpack query params
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
		state := params.Get("state")
		if state != c.State {
			glbLogger.Println("State mismatch error")
			ch <- -1
		}

		authCode := params.Get("code")
		if authCode == "" {
			glbLogger.Println("No auth token: ", authCode)
			ch <- -1
		}

		// Get access_token and refresh_token
		_, err = c.getAccessToken(authCode)
		if err != nil {
			glbLogger.Println("Error getting access token: ", err)
			ch <- -1
		}
		ch <- 0
	}

	// route handlers
	http.HandleFunc("/callback", handleCallback)

	return http.ListenAndServe(server_url, nil)
}

func (c *Credentials) getAccessToken(authCode string) (string, error) {
	body := url.Values{}
	body.Add("code", authCode)
	body.Add("redirect_uri", "http://"+redirect_url)
	body.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", access_tok_url, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}

	authHeader := c.getAuthorizationHeader()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authHeader)

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resp.Status, errors.New("error fetching access token")
	}

	authResponse := Auth{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&authResponse)
	if err != nil {
		return "", err
	}

	// Set access and refresh_token to viper config
	viper.Set("spfy-access", authResponse.AccessToken)
	viper.Set("spfy-refresh", authResponse.RefreshToken)
	err = viper.WriteConfig()
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

func (c *Credentials) getAuthorizationHeader() string {
	data := fmt.Appendf([]byte{}, "%s:%s", c.ClientId, c.ClientSecret)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))

	base64.StdEncoding.Encode(dst, data)

	return fmt.Sprintf("Basic %s", string(dst))
}
