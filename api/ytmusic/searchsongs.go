package ytmusic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func SearchSongYT(songName []string) ([]string, error) {
	// Init context and req client
	ctx := initContext()
	client := &http.Client{}

	// Init a search body
	body := SearchRequestBody{
		Ctx: ctx,
		// Query:  songName,
		Params: PARAM,
	}

	// Init return struct and videoIds array
	var videoIds []string
	var wg sync.WaitGroup
	ch := make(chan string, len(songName))

	for _, sn := range songName {
		body.Query = sn
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		go getSongInfo(*client, reqBody, &wg, ch)
	}

	wg.Wait()
	close(ch)

	for i := range ch {
		videoIds = append(videoIds, i)
	}

	return videoIds, nil
}

func getSongInfo(client http.Client, reqBody []byte, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()

	ret := ResponseStruct{}

	req, err := http.NewRequest("POST", YTMUSIC_SEARCH, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	err = json.Unmarshal(respBody, &ret)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	vidId := getVideoId(&ret)
	ch <- vidId
}

func getVideoId(ret *ResponseStruct) string {
	// I hate how google structured their data here :(
	// The following code will be hard to look at.
	if ret == nil {
		return ""
	}

	tab := ret.Contents.TabbedSearchResultsRenderer.Tabs
	if len(tab) == 0 {
		return ""
	}

	sectionListContents := tab[0].TabRenderer.Content.SectionListRenderer.Contents
	if len(sectionListContents) == 0 {
		return ""
	}

	musicShelfContent := sectionListContents[0].MusicShelfRenderer.Contents
	if len(musicShelfContent) == 0 {
		musicShelfContent = sectionListContents[1].MusicShelfRenderer.Contents
		// If the song is not in the first two suggested, return nothing.
		//	Not worth checking the others.
		if len(musicShelfContent) == 0 {
			return ""
		}
	}

	return musicShelfContent[0].MusicResponsiveListItemRenderer.PlaylistItemData.VideoID
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
