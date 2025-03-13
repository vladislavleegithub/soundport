package ytmusic

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func getAuthToken(cookie string) (string, error) {
	// Check config for auth token
	authToken := viper.GetString("yt-auth-token")
	if len(authToken) > 0 {
		return authToken, nil
	}

	// If config doesnt have auth token
	spsid, err := getSpsidFromCookie(cookie)
	if err != nil {
		return "", err
	}

	authToken = constructAuthToken(spsid)

	// Add auth token to config
	viper.Set("yt-auth-token", authToken)
	err = viper.WriteConfig()
	if err != nil {
		return "", errors.New("unable to write auth header to config: " + err.Error())
	}

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
