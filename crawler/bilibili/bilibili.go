package bilibili

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	body := &body{}
	if err := json.Unmarshal(data, body); err != nil {
		return nil, err
	} else if body.Code != 0 {
		return nil, fmt.Errorf("body code: %d", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, keyword := range body.Data.Trending.List {
		title := strings.TrimSpace(keyword.ShowName)
		board.Append1(title)
	}
	return board, nil
}

type body struct {
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
