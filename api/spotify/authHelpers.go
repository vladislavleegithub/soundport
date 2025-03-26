package spotify

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// User specific url where they accept or reject login request
func getAuthURL(clientId, state string) string {
	u, err := url.Parse(auth_url)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", clientId)
	q.Set("scope", scope)
	q.Set("redirect_uri", "http://"+redirect_url)
	q.Set("state", state)

	u.RawQuery = q.Encode()

	return u.String()
}

// Generates a "state" string required by spotify auth api
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func updateSpotifyConfig(resp *authResponse) error {
	viper.Set("spfy-access", resp.AccessToken)

	if resp.RefreshToken != "" {
		viper.Set("spfy-refresh", resp.RefreshToken)
	}

	// Add token expiry
	expiresAt := time.Now().Add(resp.ExpiresIn * time.Second)
	viper.Set("spfy-expires-at", expiresAt)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func makeAuthRequest(body *strings.Reader) error {
	req, err := http.NewRequest("POST", token_url, body)
	if err != nil {
		return err
	}

	authHeader := getAuthorizationHeader()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authHeader)

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
