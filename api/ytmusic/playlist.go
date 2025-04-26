package ytmusic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type CreatePlaylist struct {
	Ctx           *Context   `json:"context"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	PrivacyStatus StatusType `json:"privacyStatus"`
}

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

	bodyBytes, _ := io.ReadAll(resp.Body)
	if len(bodyBytes) == 0 {
		return "", fmt.Errorf("empty response body")
	}

	err = json.Unmarshal(bodyBytes, &respStruct)
	if err != nil {
		return "", fmt.Errorf("error decoding playlistid: %s", err)
	}

	if respStruct.PlaylistID == "" {
		return "", fmt.Errorf("unable to create playlist. Something went wrong")
	}

	return respStruct.PlaylistID, nil
}
