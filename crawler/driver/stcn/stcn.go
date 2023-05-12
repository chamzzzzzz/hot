package stcn

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "stcn"
	ProxySwitch = false
	URL         = "https://www.stcn.com/article/index-hot-list.html"
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
	for _, data := range body.Data {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: "https://www.stcn.com" + strings.TrimSpace(data.URL)})
	}
	return board, nil
}

type body struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
	Data  []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"data"`
}

func (body *body) NormalizedCode() int {
	if body.State == 1 {
		return 0
	} else {
		return 1
	}
}
