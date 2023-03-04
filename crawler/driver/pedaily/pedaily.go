package pedaily

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "pedaily"
	ProxySwitch = false
	URL         = "https://www.pedaily.cn"
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
	for _, ul := range dom.QueryAll("ul", "class", "news-hot-list") {
		for _, a := range ul.QueryAll("a") {
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Href())
			p, err := a.Find("p")
			if err == nil {
				title = strings.TrimSpace(p.Text())
			}
			if title != "" {
				board.AppendTitleURL(title, url)
			}
		}
	}
	return board, nil
}
