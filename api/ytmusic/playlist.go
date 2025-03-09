package ytmusic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type StatusType string

func PlaylistAdd(name string, status StatusType, songs []string) error {
	// init context
	ctx := newContext()

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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))
	return nil
}
