package ytmusic

import (
	"bytes"
	"encoding/json"
	"sync"
)

const batchSize = 50

type findTracksReq struct {
	Ctx    *Context `json:"context"`
	Query  string   `json:"query"`
	Params string   `json:"params"`
}

type actions struct {
	VideoId      string `json:"addedVideoId"`
	Action       string `json:"action"`
	DeDupeOption string `json:"dedupeOption"`
}

type addTracks struct {
	Ctx        *Context  `json:"context"`
	PlaylistID string    `json:"playlistId"`
	Actions    []actions `json:"actions"`
}

func (c *Client) AddTracks(plId string, tracks []string) (int, bool) {
	total_songs_added := 0

	for start, end := 0, 0; start <= len(tracks)-1; start = end {
		end = min(start+batchSize, len(tracks))

		batch := tracks[start:end]
		vidIdChan := make(chan string, len(batch))

		c.findTracks(batch, vidIdChan)

		vidIdList := []actions{}
		// Construct the 'actions' field required by ytmusic api
		for vidId := range vidIdChan {
			if vidId != "" {
				vidIdList = append(vidIdList, actions{
					VideoId:      vidId,
					Action:       "ACTION_ADD_VIDEO",
					DeDupeOption: "DEDUPE_OPTION_SKIP",
				})
				total_songs_added += 1
			}
		}
		body := addTracks{
			Ctx:        c.ctx,
			Actions:    vidIdList,
			PlaylistID: plId,
		}

		reqBody, err := json.Marshal(body)
		if err != nil {
			glbLogger.Println("Error constructing body: ", err)
			return 0, false
		}

		_, err = c.makeRequest(YTMUSIC_PLAYLIST_UPDATE, bytes.NewBuffer(reqBody))
		if err != nil {
			glbLogger.Println("Error sending request: ", err)
			return 0, false
		}

	}

	return total_songs_added, true
}

func (c *Client) findTracks(songs []string, ch chan<- string) {
	var wg sync.WaitGroup

	for _, song := range songs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			body := findTracksReq{
				Ctx:    c.ctx,
				Params: PARAM,
				Query:  song,
			}

			reqBody, err := json.Marshal(body)
			if err != nil {
				glbLogger.Println("Error constructing body: ", err)
				return
			}
			resp, err := c.makeRequest(YTMUSIC_SEARCH, bytes.NewBuffer(reqBody))
			if err != nil {
				glbLogger.Println("Error constructing body: ", err)
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
			if vidId == "" {
				nfLogger.Println("Not found: ", song)
			}
			ch <- vidId
		}()
	}
	wg.Wait()
	close(ch)
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
