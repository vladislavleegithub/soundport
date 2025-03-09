package ytmusic

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func header(req *http.Request) error {
	// Init base headers
	initHeaders(req)

	visitorId, err := GetVisitorId()
	if err != nil {
		fmt.Println("error getting visitor id: ", err)
		return err
	}

	cookie := viper.GetString("yt-cookie")
	if len(cookie) == 0 {
		// FIX: Fix error message
		return errors.New("unset cookie. Please set")
	}

	authHeader, err := getAuthHeader(cookie)
	if err != nil {
		fmt.Println("error getting auth header: ", err)
		return err
	}

	// Add the remaining two headers
	req.Header.Add("X-Goog-Visitor-Id", visitorId)
	req.Header.Add("authorization", authHeader)
	req.Header.Add("Cookie", cookie)

	return nil
}

func newContext() *Context {
	context := &Context{}
	context.Client.Hl = "en"
	context.Client.Gl = "IN"
	context.Client.ClientName = "WEB_REMIX"
	context.Client.ClientVersion = "1." + time.Now().Format("20060102") + ".01.00"

	return context
}

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

func getAuthHeader(cookie string) (string, error) {
	var authHeader string

	// Check config for auth
	authHeader = viper.GetString("authorization")
	if len(authHeader) > 0 {
		return authHeader, nil
	}

	spsid, err := getSpsidFromCookie(cookie)
	if err != nil {
		return "", err
	}

	authHeader = constructAuthToken(spsid)

	// Add authHeader to config
	viper.Set("authorization", authHeader)
	err = viper.WriteConfig()
	if err != nil {
		return "", errors.New("unable to write auth header to config: " + err.Error())
	}

	return authHeader, nil
}

// Some voodoo stuff, reverse engineered by
// https://stackoverflow.com/a/32065323/5726546
func constructAuthToken(spsid string) string {
	unixTimeStamp := strconv.Itoa(int(time.Now().Unix()))
	auth := spsid + " " + YTMUSIC_BASE_URL

	hasher := sha1.New()
	hasher.Write([]byte(unixTimeStamp + " " + auth))

	sha := hex.EncodeToString(hasher.Sum(nil))
	authHeader := "SAPISIDHASH " + unixTimeStamp + "_" + sha

	return authHeader
}

func getSpsidFromCookie(cookie string) (string, error) {
	paresdCookie, err := http.ParseCookie(cookie)
	if err != nil {
		return "", err
	}

	var spsid string
	for _, v := range paresdCookie {
		if v.Name == "__Secure-3PAPISID" {
			spsid = v.Value
			break
		}
	}

	return spsid, nil
}
