package ytmusic

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
)

func GetVisitorId() (string, error) {
	// Get response
	response, err := sendGetRequest()
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.New("Unable to fetch data: " + response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// extract visitor Id
	regex := `ytcfg\.set\s*\(\s*({.+?})\s*\)\s*;`
	matches, err := extractData(regex, body)
	if err != nil {
		return "", err
	}

	visitorIdStruct := struct {
		VisitorId string `json:"VISITOR_DATA"`
	}{}
	err = json.Unmarshal([]byte(matches[1]), &visitorIdStruct)
	if err != nil {
		return "", err
	}

	return visitorIdStruct.VisitorId, nil
}

func extractData(regex string, body []byte) ([]string, error) {
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	matches := r.FindStringSubmatch(string(body))

	return matches, nil
}
