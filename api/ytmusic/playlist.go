package ytmusic

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c *Client) CreatePlaylist(name string, desc string) (string, error) {
	body := CreatePlaylist{
		Ctx:           c.ctx,
		Title:         name,
		PrivacyStatus: PRIVATE,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	resp, err := c.makeRequest(YTMUSIC_PLAYLIST, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respStruct := struct {
		PlaylistID string `json:"playlistId"`
	}{}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&respStruct)
	if err != nil {
		return "", err
	}

	if respStruct.PlaylistID == "" {
		return "", fmt.Errorf("unable to create playlist. Something went wrong")
	}

	return respStruct.PlaylistID, nil
}
