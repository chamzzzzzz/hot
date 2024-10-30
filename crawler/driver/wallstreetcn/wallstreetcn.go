package wallstreetcn

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "wallstreetcn"
	ProxySwitch = false
	URL         = "https://api-one-wscn.awtmt.com/apiv1/content/articles/hot?period=all"
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
	for _, data := range body.Data.DayItems {
		board.Append(&hot.Hot{Title: strings.TrimSpace(data.Title), URL: strings.TrimSpace(data.URI)})
	}
	return board, nil
}

type body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		DayItems []struct {
			Title string `json:"title"`
			URI   string `json:"uri"`
		} `json:"day_items"`
	} `json:"data"`
}

func (body *body) NormalizedCode() int {
	if body.Code == 20000 {
		return 0
	} else {
		return 1
	}
}
