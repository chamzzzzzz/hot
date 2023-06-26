package voachinese

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "voachinese"
	ProxySwitch = true
	URL         = "https://www.voachinese.com"
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
	div, err := dom.Find("div", "class", "row trends-wg")
	if err != nil {
		return nil, err
	}
	for _, a := range div.QueryAll("a") {
		h4, err := a.Find("h4", "class", "trends-wg__item-txt")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(h4.Text())
		url := "https://www.voachinese.com" + strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{Title: title, URL: url})
	}
	return board, nil
}
