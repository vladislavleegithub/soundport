package ytmusic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// FIX: ADD PROPER ERROR RETURNS
func SearchSongYT(songName []string) ([]string, error) {
	// Init context and req client
	ctx := newContext()
	client := &http.Client{}

	// Init a search body
	body := SearchRequestBody{
		Ctx: ctx,
		// Query:  songName,
		Params: PARAM,
	}

	// Init return struct and videoIds array
	ret := ResponseStruct{}
	var videoIds []string

	for _, sn := range songName {
		body.Query = sn
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", YTMUSIC_SEARCH, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(respBody, &ret)
		if err != nil {
			return nil, err
		}

		vidId := getVideoId(&ret)
		videoIds = append(videoIds, vidId)
	}

	return videoIds, nil
}

func getVideoId(ret *ResponseStruct) string {
	return ret.Contents.TabbedSearchResultsRenderer.Tabs[0].TabRenderer.Content.SectionListRenderer.Contents[0].MusicShelfRenderer.Contents[0].MusicResponsiveListItemRenderer.PlaylistItemData.VideoID
}

type ResponseStruct struct {
	Contents struct {
		TabbedSearchResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents []struct {
								MusicShelfRenderer struct {
									Title struct {
										Runs []struct {
											Text string `json:"text"`
										} `json:"runs"`
									} `json:"title"`
									Contents []struct {
										MusicResponsiveListItemRenderer struct {
											PlaylistItemData struct {
												VideoID string `json:"videoId"`
											} `json:"playlistItemData"`
										} `json:"musicResponsiveListItemRenderer"`
									} `json:"contents"`
								} `json:"musicShelfRenderer"`
							} `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"tabbedSearchResultsRenderer"`
	} `json:"contents"`
}
