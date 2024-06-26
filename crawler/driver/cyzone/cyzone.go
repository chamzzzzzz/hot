package cyzone

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "cyzone"
	ProxySwitch = false
	URL         = "https://www.cyzone.cn"
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
	for _, div := range dom.QueryAll("div", "class", "side-panel") {
		span := div.Query("div", "class", "panel-tit").Query("span")
		if span == nil {
			continue
		}
		if span.Text() != "热门文章" {
			continue
		}
		for _, a := range div.QueryAll("div", "class", "item").Query("div", "class", "intro").Query("a") {
			title := strings.TrimSpace(a.Text())
			url := "https://www.cyzone.cn" + strings.TrimSpace(a.Href())
			board.Append(&hot.Hot{Title: title, URL: url})
		}
	}
	return board, nil
}
