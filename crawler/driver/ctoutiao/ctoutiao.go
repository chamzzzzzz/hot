package ctoutiao

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "ctoutiao"
	ProxySwitch = false
	URL         = "https://www.ctoutiao.com"
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
	for _, div := range dom.QueryAll("div", "class", "ctt-hottext-item") {
		a, err := div.Find("a")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Text())
		url := "https://www.ctoutiao.com" + strings.TrimSpace(a.Href())
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
