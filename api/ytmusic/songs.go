package ytmusic

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"
)

func SearchSong(songName []string) ([]string, error) {
	// Init context
	ctx := initContext()
	client := &http.Client{}

	body := SearchRequestBody{
		Ctx:    ctx,
		Params: PARAM,
	}

	// Returned video Ids
	var videoIds []string
	var wg sync.WaitGroup
	ch := make(chan string, len(songName))

	// Run a go routine for every song name
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
		if i == "" {
			// TODO: Handle not found songs.
			// Currently we are just skipping them without informing.
			continue
		}
		videoIds = append(videoIds, i)
	}

	return videoIds, nil
}

func getSongInfo(client http.Client, reqBody []byte, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()

	req, err := http.NewRequest("POST", YTMUSIC_SEARCH, bytes.NewBuffer(reqBody))
	if err != nil {
		glbLogger.Println("Error constructing request: ", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		glbLogger.Println("Error sending request: ", err)
		return
	}
	defer resp.Body.Close()

	// Read body into a struct
	ret := ResponseStruct{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ret)
	if err != nil {
		glbLogger.Println("Error reading response body: ", err)
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
