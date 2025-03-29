package ytmusic

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/Samarthbhat52/soundport/api/types"
)

const batchSize = 50

func (c *Client) FindTracks(songs []string, ch chan<- types.SongDetails) {
	for start, end := 0, 0; start <= len(songs)-1; start = end {
		end = min(start+batchSize, len(songs))

		batchedSongs := songs[start:end]
		c.batchProcess(batchedSongs, ch)
	}

	close(ch)
}

func (c *Client) batchProcess(songs []string, ch chan<- types.SongDetails) {
	var wg sync.WaitGroup
	for _, song := range songs {
		wg.Add(1)

		body := SearchRequestBody{
			Ctx:    c.ctx,
			Params: PARAM,
			Query:  song,
		}
		reqBody, err := json.Marshal(body)
		if err != nil {
			glbLogger.Println("Error constructing body: ", err)
			return
		}

		req, err := http.NewRequest("POST", YTMUSIC_SEARCH, bytes.NewBuffer(reqBody))
		if err != nil {
			glbLogger.Println("Error constructing request: ", err)
			return
		}

		go func() {
			defer wg.Done()

			resp, err := c.client.Do(req)
			if err != nil {
				glbLogger.Println("Error sending request: ", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				glbLogger.Println("Error sending request: ", resp.StatusCode)
				return
			}

			// Read body into a struct
			ret := ResponseStruct{}
			decoder := json.NewDecoder(resp.Body)
			err = decoder.Decode(&ret)
			if err != nil {
				glbLogger.Println("Error reading response body: ", err)
				return
			}

			vidId := getVideoId(&ret)
			if vidId != "" {
				ch <- types.SongDetails{
					Name:  song,
					Id:    vidId,
					Found: true,
				}
			} else {
				ch <- types.SongDetails{
					Name:  song,
					Id:    "NULL",
					Found: false,
				}
			}
		}()
	}
	wg.Wait()
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
