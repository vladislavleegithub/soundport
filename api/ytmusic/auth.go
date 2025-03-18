package ytmusic

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

func GetAuthToken(cookie string) (string, error) {
	// If config doesnt have auth token
	spsid, err := getSpsidFromCookie(cookie)
	if err != nil {
		return "", err
	}
	authToken := constructAuthToken(spsid)

	return authToken, nil
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
