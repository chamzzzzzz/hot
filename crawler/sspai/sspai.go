package sspai

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
	return "sspai"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sspai.com/api/v1/article/tag/page/get?limit=10&offset=0&tag=热门文章&released=false", nil)
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
	} else if bodyJson.Error != 0 {
		return nil, fmt.Errorf("body error: %d", bodyJson.Error)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson.Data {
		board.Append(fmt.Sprintf("%s | %s", data.Title, data.Summary), fmt.Sprintf("/post/%d", data.ID), date)
	}
	return board, nil
}

type bodyJson struct {
	Error int    `json:"error"`
	Msg   string `json:"msg"`
	Data  []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
	} `json:"data"`
}
