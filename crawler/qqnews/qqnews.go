package qqnews

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
	return "qqnews"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", `https://i.news.qq.com/trpc.qqnews_web.kv_srv.kv_srv_http_proxy/list?sub_srv_id=24hours&srv_id=pc&offset=0&limit=20&strategy=1&ext={%22pool%22:[%22top%22],%22is_filter%22:7,%22check_type%22:true}`, nil)
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
	} else if body.Ret != 0 {
		return nil, fmt.Errorf("body ret: %d", body.Ret)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data.List {
		date, err := time.ParseInLocation("2006-01-02 15:04:05", data.PublishTime, time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(data.Title, "", data.URL, date)
	}
	return board, nil
}

type body struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		List []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			PublishTime string `json:"publish_time"`
		} `json:"list"`
	} `json:"data"`
}
