package bbc

import (
	"encoding/json"
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	body := &body{}
	if err := json.Unmarshal(data, body); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, record := range body.Records {
		title := strings.TrimSpace(record.Promo.Headlines.ShortHeadline)
		summary := strings.TrimSpace(record.Promo.Summary)
		url := "https://www.bbc.com/" + strings.TrimSpace(strings.Trim(record.Promo.ID, "urn:bbc:ares::asset:"))
		date := time.UnixMilli(record.Promo.Timestamp)
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	Records []struct {
		Promo struct {
			Headlines struct {
				ShortHeadline string `json:"shortHeadline"`
				Headline      string `json:"headline"`
			} `json:"headlines"`
			Summary   string `json:"summary"`
			Timestamp int64  `json:"timestamp"`
			ID        string `json:"id"`
		} `json:"promo,omitempty"`
	} `json:"records"`
}
