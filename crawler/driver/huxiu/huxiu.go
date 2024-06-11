package huxiu

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "huxiu"
	ProxySwitch = false
	URL         = "https://www.huxiu.com"
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
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, a := range dom.Query("div", "class", "hot-article-wrap").QueryAll("a", "class", "title-wrap two-lines") {
		h3, err := a.Find("h3")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(h3.Text())
		url := strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{Title: title, URL: url})
	}
	return board, nil
}
