package ytmusic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

func makeRequest(client http.Client, reqBody []byte, wg *sync.WaitGroup, ch chan<- string) {
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
