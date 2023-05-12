package v2ex

import (
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "v2ex"
	ProxySwitch = true
	URL         = "https://www.v2ex.com/api/topics/hot.json"
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
	var body []body
	if err := httputil.Request("GET", URL, nil, "json", &body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, item := range body {
		date := time.Unix(item.Created, 0)
		board.Append(&hot.Hot{Title: item.Title, URL: item.URL, PublishDate: date})
	}
	return board, nil
}

type body struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Created int64  `json:"created"`
}
