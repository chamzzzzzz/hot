package sspai

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	DriverName = "sspai"
)

type Driver struct {
}

func (driver *Driver) Open(option driver.Option) (driver.Crawler, error) {
	return &Crawler{Option: option}, nil
}

func init() {
	driver.Register(DriverName, &Driver{})
}

type Crawler struct {
	Option driver.Option
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	body := &body{}
	if err := json.Unmarshal(data, body); err != nil {
		return nil, err
	} else if body.Error != 0 {
		return nil, fmt.Errorf("body error: %d", body.Error)
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		title := strings.TrimSpace(data.Title)
		summary := strings.TrimSpace(data.Summary)
		url := fmt.Sprintf("https://sspai.com/post/%d", data.ID)
		date := time.Unix(data.ReleasedTime, 0)
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	Error int    `json:"error"`
	Msg   string `json:"msg"`
	Data  []struct {
		ID           int    `json:"id"`
		Title        string `json:"title"`
		Summary      string `json:"summary"`
		ReleasedTime int64  `json:"released_time"`
	} `json:"data"`
}
