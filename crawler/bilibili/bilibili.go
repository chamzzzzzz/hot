package bilibili

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "bilibili"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.bilibili.com/x/web-interface/search/square?limit=10", nil)
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
	} else if bodyJson.Code != 0 {
		return nil, fmt.Errorf("body code: %d", bodyJson.Code)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, keyword := range bodyJson.Data.Trending.List {
		board.Append(keyword.ShowName, "", date)
	}
	return board, nil
}

type bodyJson struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	TTL     int    `json:"ttl,omitempty"`
	Data    struct {
		Trending struct {
			Title   string `json:"title,omitempty"`
			Trackid string `json:"trackid,omitempty"`
			List    []struct {
				Keyword  string `json:"keyword,omitempty"`
				ShowName string `json:"show_name,omitempty"`
				Icon     string `json:"icon,omitempty"`
				URI      string `json:"uri,omitempty"`
				Goto     string `json:"goto,omitempty"`
			} `json:"list,omitempty"`
		} `json:"trending,omitempty"`
	} `json:"data,omitempty"`
}
