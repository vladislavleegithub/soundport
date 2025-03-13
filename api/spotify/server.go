package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func (c *Credentials) StartHttpServer(ch chan string) error {
	handleCallback := func(w http.ResponseWriter, r *http.Request) {
		// Unpack query params
		params, err := url.ParseQuery(r.URL.String())
		if err != nil {
			fmt.Println("Error decoding query params: ", err.Error())
		}

		error := params.Get("error")
		if error != "" {
			fmt.Println("Error authenticating: ", error)
			ch <- error
			return
		}

		// TODO: Compare state token to the one created serverside. later.
		authCode := params.Get("/callback?code")
		if authCode == "" {
			fmt.Println("No auth token: ", authCode)
			ch <- "error"
			return
		}

		// Get access_token and refresh_token
		status, err := c.getAccessToken(authCode)
		if err != nil {
			fmt.Println("ERROR GETTING ACCESS TOKEN: ", err)
		}
		ch <- status
	}

	log.Infof("Listening to server on address: %s", server_url)

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
	fmt.Println("BODY: ", authResponse)
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
