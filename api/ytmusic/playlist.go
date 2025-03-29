package ytmusic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type StatusType string

func (c *Client) AddPlaylist(title string, songs []string) error {
	glbLogger.Printf("TOTAL SONGS FOUND: %d", len(songs))

	// prep body
	body := CreatePlaylistRequestBody{
		Ctx:           c.ctx,
		Title:         title,
		PrivacyStatus: "PRIVATE",
		VideoIds:      songs,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", YTMUSIC_PLAYLIST, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	// Prep headers
	err = constructHeader(req)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got error code: %d", resp.StatusCode)
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func constructHeader(req *http.Request) error {
	// Init base headers
	initHeaders(req)

	visitorId, err := getVisitorId()
	if err != nil {
		glbLogger.Println("error getting visitor id: ", err)
		return err
	}

	cookie := viper.GetString("yt-cookie")
	authHeader := viper.GetString("yt-auth-token")

	// Add the remaining two headers
	req.Header.Add("X-Goog-Visitor-Id", visitorId)
	req.Header.Add("authorization", authHeader)
	req.Header.Add("Cookie", cookie)

	return nil
}
