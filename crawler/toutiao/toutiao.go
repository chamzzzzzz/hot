package toutiao

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
	return "toutiao"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://i-lq.snssdk.com/api/feed/hotboard_online/v1/?category=hotboard_online&count=50", nil)
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
	} else if bodyJson.Message != "success" {
		return nil, fmt.Errorf("body message: %s", bodyJson.Message)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson.Data {
		contentJson := &contentJson{}
		if err := json.Unmarshal([]byte(data.Content), contentJson); err != nil {
			return nil, err
		}
		board.Append(contentJson.RawData.Title, "", date)
	}
	return board, nil
}

type bodyJson struct {
	Message string `json:"message"`
	Data    []struct {
		Content string `json:"content"`
		Code    string `json:"code"`
	} `json:"data"`
}

type contentJson struct {
	RawData struct {
		Title string `json:"title"`
	} `json:"raw_data"`
}
