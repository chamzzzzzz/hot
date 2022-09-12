package futu

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
	return "futu"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.futunn.com/search-stock/hot-news", nil)
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
	} else if bodyJson.Code != "0" {
		return nil, fmt.Errorf("body code: %s", bodyJson.Code)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson.Data {
		board.Append(data.Title, data.URL, date)
	}
	return board, nil
}

type bodyJson struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		PostID int    `json:"post_id"`
		Title  string `json:"title"`
		URL    string `json:"url"`
	} `json:"data"`
}
