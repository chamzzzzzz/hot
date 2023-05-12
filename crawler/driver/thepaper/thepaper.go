package thepaper

import (
	"fmt"
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "thepaper"
	ProxySwitch = false
	URL         = "https://www.thepaper.cn/contentapi/wwwIndex/rightSidebar"
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
	if body.ResultCode != 1 {
		return nil, fmt.Errorf("body result code: %d", body.ResultCode)
	}

	board := hot.NewBoard(c.Name())
	for _, news := range body.Data.HotNews {
		title := strings.TrimSpace(news.Name)
		url := "https://www.thepaper.cn/newsDetail_forward_" + strings.TrimSpace(news.ContID)
		date := time.UnixMilli(news.PubTimeLong)
		board.Append(&hot.Hot{Title: title, URL: url, PublishDate: date})
	}
	return board, nil
}

type body struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	Data       struct {
		HotNews []struct {
			ContID      string `json:"contId"`
			Name        string `json:"name"`
			PubTimeLong int64  `json:"pubTimeLong"`
		} `json:"hotNews"`
	} `json:"data"`
}
