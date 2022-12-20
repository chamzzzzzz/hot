package thehackernews

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "thehackernews"
	ProxySwitch = false
	URL         = "https://thehackernews.com"
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
	div := dom.FindStrict("div", "id", "PopularPosts1")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a", "class", "pop-link cf") {
		div := a.FindStrict("div", "class", "pop-title")
		if div.Error != nil {
			return nil, div.Error
		}
		title := strings.TrimSpace(div.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
