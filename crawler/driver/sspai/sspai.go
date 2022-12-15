package sspai

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "sspai"
	ProxySwitch = false
	URL         = "https://sspai.com/api/v1/article/tag/page/get?limit=10&offset=0&tag=热门文章&released=false"
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
	body := &body{}
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}
	if body.Error != 0 {
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
