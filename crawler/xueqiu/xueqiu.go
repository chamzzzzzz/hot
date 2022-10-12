package xueqiu

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
)

type Crawler struct {
	Cookie string
}

func (c *Crawler) Name() string {
	return "xueqiu"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	err := c.updateCookie()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://xueqiu.com/query/v1/status/hots.json?count=10&page=1&scope=day&type=news", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.Set("Cookie", c.Cookie)

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
	} else if body.Code != 200 {
		return nil, fmt.Errorf("body code: %d", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		title := strings.TrimSpace(data.Title)
		url := "https://xueqiu.com" + strings.TrimSpace(data.Target)
		board.AppendTitleURL(title, url)
	}
	return board, nil
}

func (c *Crawler) updateCookie() error {
	if c.Cookie != "" {
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("HEAD", "https://xueqiu.com", nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	for _, cookie := range res.Cookies() {
		req.AddCookie(cookie)
	}
	c.Cookie = req.Header.Get("Cookie")
	return nil
}

type body struct {
	Code int `json:"code"`
	Data []struct {
		Target string `json:"target"`
		Text   string `json:"text"`
		Title  string `json:"title"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}
