package ytmusic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FIX: ADD PROPER ERROR RETURNS
func SearchSongYT(songName string) {
	// Init context
	ctx := newContext()

	// Init a search body
	body := SearchRequestBody{
		BaseRequestBody: BaseRequestBody{
			Ctx: ctx,
		},
		Query:  songName,
		Params: PARAM,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("json marshal error: ", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", YTMUSIC_SEARCH, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("req prep error: ", err)
		return
	}

	initHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error fetching song: ", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading resp body: ", err)
		return
	}

	ret := ResponseStruct{}
	err = json.Unmarshal(respBody, &ret)
	if err != nil {
		fmt.Println("error unmarshalling resp: ", err)
		return
	}

	vidId := getVideoId(&ret)
	fmt.Println(vidId)
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
