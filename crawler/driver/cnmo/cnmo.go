package cnmo

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "cnmo"
	ProxySwitch = false
	URL         = "https://www.cnmo.com"
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
	for _, ul := range dom.QueryAll("ul", "class", "hot-7") {
		if ul.Attribute("style") != "display: block;" {
			continue
		}
		for _, a := range ul.QueryAll("div", "class", "hot-7-tit").Query("a") {
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Href())
			if strings.HasPrefix(url, "//") {
				url = "https:" + url
			}
			board.Append(&hot.Hot{Title: title, URL: url})
		}
	}
	return board, nil
}
