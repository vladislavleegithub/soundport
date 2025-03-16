package spotify

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Samarthbhat52/soundport/logger"
)

var glbLogger = logger.GetInstance()

func (c *credentials) StartHttpServer(ch chan int, state string) error {
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
		retState := params.Get("state")
		if retState != state {
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

func (c *credentials) getAccessToken(authCode string) (string, error) {
	body := url.Values{}
	body.Add("code", authCode)
	body.Add("redirect_uri", "http://"+redirect_url)
	body.Add("grant_type", "authorization_code")

	req, err := http.NewRequest("POST", token_url, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}

	authHeader := c.getAuthorizationHeader()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", authHeader)

	err = handleAuthResponse(req)
	if err != nil {
		return "", err
	}

	return "processed", nil
}

func (c *credentials) getAuthorizationHeader() string {
	data := fmt.Appendf([]byte{}, "%s:%s", c.clientId, c.clientSecret)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))

	base64.StdEncoding.Encode(dst, data)

	return fmt.Sprintf("Basic %s", string(dst))
}

func (c *credentials) GetAuthURL(state string) string {
	u, err := url.Parse(auth_url)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Set("response_type", "code")
	q.Set("client_id", c.clientId)
	q.Set("scope", scope)
	q.Set("redirect_uri", "http://"+redirect_url)
	q.Set("state", state)

	// Encode the Query
	u.RawQuery = q.Encode()

	return u.String()
}
