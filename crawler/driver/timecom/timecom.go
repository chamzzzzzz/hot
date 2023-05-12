package timecom

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "timecom"
	ProxySwitch = true
	URL         = "https://time.com/"
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
	for _, li := range dom.QueryAll("li", "class", "most-popular-feed__item") {
		for i, a := range li.QueryAll("a") {
			if i == 1 {
				h3, err := a.Find("h3")
				if err != nil {
					return nil, err
				}
				title := strings.TrimSpace(h3.Text())
				url := strings.TrimSpace(a.Href())
				board.Append(&hot.Hot{Title: title, URL: url})
			}
		}
	}
	return board, nil
}
