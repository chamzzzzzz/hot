package v2ex

import (
	"encoding/json"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Crawler struct {
	Proxy string
}

func (c *Crawler) Name() string {
	return "v2ex"
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

	req, err := http.NewRequest("GET", "https://www.v2ex.com/api/topics/hot.json", nil)
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

	var bodyJson bodyJson
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, item := range bodyJson {
		board.Append(item.Title, item.URL, date)
	}
	return board, nil
}

type bodyJson []struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}
