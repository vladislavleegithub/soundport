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

func handleAuthResponse(req *http.Request) error {
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
