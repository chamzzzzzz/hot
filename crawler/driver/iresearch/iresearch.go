package iresearch

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "iresearch"
	ProxySwitch = false
	URL         = "https://www.iresearch.cn"
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
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.DetectContentEncoding = true
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div, err := dom.Find("div", "class", "g-news-column")
	if err != nil {
		return nil, err
	}
	for _, a := range div.QueryAll("a") {
		title := strings.TrimSpace(a.FullText())
		url := strings.TrimSpace(a.Href())
		if title == "" {
			continue
		}
		board.Append(&hot.Hot{Title: title, URL: url})
	}
	return board, nil
}
