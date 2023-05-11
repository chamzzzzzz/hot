package gizmodo

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "gizmodo"
	ProxySwitch = false
	URL         = "https://gizmodo.com/"
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
	for _, li := range dom.QueryAll("li", "class", "sc-4yr8xe-2 sc-41vfnj-1") {
		a, err := li.Find("a")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.FullText())
		url := strings.TrimSpace(a.Href())
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
