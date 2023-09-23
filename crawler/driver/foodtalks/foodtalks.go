package foodtalks

import (
	"fmt"
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "foodtalks"
	ProxySwitch = false
	URL         = "https://api-we.foodtalks.cn/news/news/hot/page?current=1&size=15"
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
	for _, data := range body.Data.Records {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: fmt.Sprintf("https://www.foodtalks.cn/news/%d", data.ID)})
	}
	return board, nil
}

type body struct {
	Code int `json:"code"`
	Data struct {
		Records []struct {
			ID    int    `json:"id"`
			Title string `json:"title"`
		} `json:"records"`
	} `json:"data"`
}

func (body *body) NormalizedCode() int {
	return body.Code
}
