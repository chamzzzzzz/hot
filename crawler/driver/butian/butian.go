package butian

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "butian"
	ProxySwitch = false
	URL         = "https://forum.butian.net"
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
	for _, div := range dom.QueryAll("div", "class", "col-md-6") {
		h4 := div.Query("h4")
		if h4 == nil {
			continue
		}
		if strings.TrimSpace(h4.Text()) != "最新文章" {
			continue
		}
		for _, a := range div.QueryAll("li").Query("a") {
			title := strings.TrimSpace(a.Title())
			url := strings.TrimSpace(a.Href())
			board.Append(&hot.Hot{Title: title, URL: url})
		}
	}
	return board, nil
}
