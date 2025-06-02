package ytmusic

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Samarthbhat52/soundport/logger"
	"github.com/spf13/viper"
)

const (
	YTMUSIC_BASE_URL        = "https://music.youtube.com"
	YTMUSIC_API             = YTMUSIC_BASE_URL + "/youtubei/v1"
	YTMUSIC_SEARCH          = YTMUSIC_API + "/search?limit=1"
	YTMUSIC_PLAYLIST        = YTMUSIC_API + "/playlist/create"
	YTMUSIC_PLAYLIST_UPDATE = YTMUSIC_API + "/browse/edit_playlist"
	PARAM                   = "EgWKAQIIAWoQEAMQBBAJEAoQBRAREBAQFQ%3D%3D"
)

type StatusType string

const (
	PRIVATE StatusType = "PRIVATE"
	PUBLIC  StatusType = "PUBLIC"
)

var (
	glbLogger = logger.GetInstance()
	nfLogger  = logger.GetNotFoundLogInstance()
)

type Context struct {
	Client struct {
		Hl            string `json:"hl"`
		Gl            string `json:"gl"`
		ClientName    string `json:"client_name"`
		ClientVersion string `json:"client_version"`
	} `json:"client"`
	User struct{} `json:"user"`
}

type Client struct {
	ctx    *Context
	header http.Header
}

func initContext() *Context {
	context := &Context{}
	context.Client.Hl = "en"
	context.Client.Gl = "IN"
	context.Client.ClientName = "WEB_REMIX"
	context.Client.ClientVersion = "1." + time.Now().Format("20060102") + ".01.00"

	return context
}

func baseHeaders() http.Header {
	header := make(http.Header)

	// Base headers
	header.Add(
		"user-agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0",
	)
	header.Add("accept", "*/*")
	header.Add("accept-language", "en-US,en;q=0.5")
	header.Add("alt-used", "music.youtube.com")
	header.Add("connection-type", "keep-alive")
	header.Add("content-type", "application/json")
	header.Add("Origin", YTMUSIC_BASE_URL)
	header.Add("X-Goog-AuthUser", "0")

	return header
}

func postHeader() (http.Header, error) {
	// Base headers
	header := baseHeaders()

	// POST auth headers
	visitorId, err := getVisitorId()
	if err != nil {
		glbLogger.Println("error getting visitor id: ", err)
		return nil, err
	}

	cookie := viper.GetString("yt-cookie")
	authHeader, err := GetAuthToken(cookie)
	if err != nil {
		glbLogger.Println("error getting auth header: ", err)
		return nil, err
	}

	// Add the remaining two headers
	header.Add("X-Goog-Visitor-Id", visitorId)
	header.Add("Authorization", authHeader)
	header.Add("Cookie", cookie)

	return header, nil
}

func NewClient() *Client {
	h, err := postHeader()
	if err != nil {
		fmt.Println("Error creating a youtube music client.")
		os.Exit(1)
	}

	return &Client{
		ctx:    initContext(),
		header: h,
	}
}

func (c *Client) makeRequest(url string, body *bytes.Buffer) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header = c.header
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got error code: %d", resp.StatusCode)
	}

	return resp, nil
}
