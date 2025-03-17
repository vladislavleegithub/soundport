package spotify

import (
	"net/http"
	"net/url"

	"github.com/Samarthbhat52/soundport/logger"
)

var glbLogger = logger.GetInstance()

func StartHttpServer(ch chan int, state string) error {
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
		_, err = getAccessToken(authCode)
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
