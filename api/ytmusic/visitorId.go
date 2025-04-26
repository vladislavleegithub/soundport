package ytmusic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func getVisitorId() (string, error) {
	// Get request
	response, err := sendGetRequest()
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return "", errors.New("Unable to fetch data: " + response.Status)
	}

	// Extract response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// extract json string that contains visitor id
	regex := `ytcfg\.set\s*\(\s*({.+?})\s*\)\s*;`
	matches, err := extractJsonFromResponse(regex, body)
	if err != nil {
		return "", err
	}

	if len(matches) < 2 {
		return "", fmt.Errorf("unable to find visitor id: %v", matches)
	}

	// Unmarshal to get visitor id
	visitorIdStruct := struct {
		VisitorId string `json:"VISITOR_DATA"`
	}{}
	err = json.Unmarshal([]byte(matches[1]), &visitorIdStruct)
	if err != nil {
		return "", err
	}

	return visitorIdStruct.VisitorId, nil
}

// Used to collect visitorId
func sendGetRequest() (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", YTMUSIC_BASE_URL, nil)
	if err != nil {
		return nil, err
	}

	headers := baseHeaders()
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error code: %v", resp.StatusCode)
	}

	return resp, nil
}

func extractJsonFromResponse(regex string, body []byte) ([]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	matches := r.FindStringSubmatch(string(body))

	return matches, nil
}
