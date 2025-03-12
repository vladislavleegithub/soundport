package ytmusic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/viper"
)

type StatusType string

func PlaylistAdd(name string, status StatusType, songs []string) error {
	// init context
	ctx := initContext()

	// prep body
	body := CreatePlaylistRequestBody{
		Ctx:           ctx,
		Title:         name,
		PrivacyStatus: status,
		VideoIds:      songs,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", YTMUSIC_PLAYLIST, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	// Prep headers
	err = header(req)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func header(req *http.Request) error {
	// Init base headers
	initHeaders(req)

	visitorId, err := GetVisitorId()
	if err != nil {
		fmt.Println("error getting visitor id: ", err)
		return err
	}

	cookie := viper.GetString("yt-cookie")
	if len(cookie) == 0 {
		// FIX: Fix error message
		return errors.New("unset cookie. Please set")
	}

	authHeader, err := getAuthToken(cookie)
	if err != nil {
		fmt.Println("error getting auth header: ", err)
		return err
	}

	// Add the remaining two headers
	req.Header.Add("X-Goog-Visitor-Id", visitorId)
	req.Header.Add("authorization", authHeader)
	req.Header.Add("Cookie", cookie)

	return nil
}
