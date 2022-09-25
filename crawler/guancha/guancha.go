package guancha

import (
	"encoding/json"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "guancha"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://user.guancha.cn/news-api/fengwen-index-list.json?page=1&order=3", nil)
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

	var bodyJson []bodyJson
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson {
		board.Append(data.Title, data.PostURL, date)
	}
	return board, nil
}

type bodyJson struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	PostURL string `json:"post_url"`
}
