package thepaper

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
	return "thepaper"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.thepaper.cn/contentapi/wwwIndex/rightSidebar", nil)
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
	} else if bodyJson.ResultCode != 1 {
		return nil, fmt.Errorf("body result code: %d", bodyJson.ResultCode)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()

	for _, news := range bodyJson.Data.HotNews {
		board.Append(news.Name, fmt.Sprintf("https://www.thepaper.cn/newsDetail_forward_%s", news.ContID), date)
	}
	return board, nil
}

type bodyJson struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	Data       struct {
		HotNews []struct {
			ContID string `json:"contId"`
			Name   string `json:"name"`
		} `json:"hotNews"`
	} `json:"data"`
}
