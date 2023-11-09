package infzm

import (
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "infzm"
	ProxySwitch = false
	URL         = "http://www.infzm.com/hot_contents?format=json"
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

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data.HotContents {
		date, err := time.ParseInLocation("2006-01-02 15:04:05", data.PublishTime, time.Local)
		if err != nil {
			return nil, err
		}
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Subject), PublishDate: date})
	}
	return board, nil
}

type body struct {
	Code int `json:"code"`
	Data struct {
		HotContents []struct {
			Subject     string `json:"subject"`
			PublishTime string `json:"publish_time"`
		} `json:"hot_contents"`
	} `json:"data"`
}

func (body *body) NormalizedCode() int {
	if body.Code == 200 {
		return 0
	}
	return body.Code
}
