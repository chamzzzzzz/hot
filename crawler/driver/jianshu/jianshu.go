package jianshu

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "jianshu"
	ProxySwitch = false
	URL         = "https://www.jianshu.com/shakespeare/v2/notes/recommend"
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
	for _, data := range body.Items {
		board.AppendTitleURL(strings.TrimSpace(data.Title), strings.TrimSpace((data.URL)))
	}
	return board, nil
}

type body struct {
	Open  int `json:"open"`
	Items []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"items"`
}

func (body *body) NormalizedCode() int {
	if body.Open == 1 {
		return 0
	} else {
		return 1
	}
}
