package bbc

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Crawler struct {
	Proxy string
}

func (c *Crawler) Name() string {
	return "bbc"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	if c.Proxy != "" {
		if proxyUrl, err := url.Parse(c.Proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	req, err := http.NewRequest("GET", "https://www.bbc.com/zhongwen/mostread/simp.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyJson := &bodyJson{}
	if err := json.Unmarshal(body, bodyJson); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, record := range bodyJson.Records {
		board.Append(fmt.Sprintf("%s | %s", record.Promo.Headlines.ShortHeadline, record.Promo.Summary), strings.Trim(record.Promo.ID, "urn:bbc:ares::asset:"), date)
	}
	return board, nil
}

type bodyJson struct {
	Records []struct {
		Promo struct {
			Headlines struct {
				ShortHeadline string `json:"shortHeadline"`
				Headline      string `json:"headline"`
			} `json:"headlines"`
			Summary string `json:"summary"`
			ID      string `json:"id"`
		} `json:"promo,omitempty"`
	} `json:"records"`
}
