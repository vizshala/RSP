package rsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	//"log"
)

// the following fields refers to bitly api documentation https://dev.bitly.com/v4_documentation.html

type bitlyShortenALinkReq struct {
	GroupGUID string `json:"group_guid"`
	Domain    string `json:"domain"`
	LongURL   string `json:"long_url"`
}

type bitlyShortenALinkRes struct {
	CreatedAt string `json:"created_at"`
	ID        string `json:"id"`
	Link      string `jons:"link"`
	// ignore other fileds
}

// CreateShortUrl takes usage of bit.ly shortening URL service
func CreateShortURL(originalUrl string) (string, int) {
	shortenLinkReq := &bitlyShortenALinkReq{
		GroupGUID: "",
		Domain:    "bit.ly",
		LongURL:   originalUrl,
	}
	reqStr, _ := json.Marshal(shortenLinkReq)

	req, err := http.NewRequest("POST", "https://api-ssl.bitly.com/v4/bitlinks", bytes.NewBuffer(reqStr))
	// the auth token is hardcoded just for demonstration purpose
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", "a05786c8cded98cabec3842df69fa999743474a0"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", 404
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	shortenLinkRes := bitlyShortenALinkRes{}

	if res.StatusCode == 200 || res.StatusCode == 201 {
		json.Unmarshal(body, &shortenLinkRes)
	}

	return shortenLinkRes.Link, res.StatusCode
}
