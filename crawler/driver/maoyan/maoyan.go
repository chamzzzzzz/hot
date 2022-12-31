package maoyan

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	News    = "news"
	Actor   = "actor"
	Unknown = "unknown"
)

const (
	DriverName  = "maoyan"
	ProxySwitch = false
	URL         = "https://www.maoyan.com"
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
	cookie string
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	if err := c.updatecookie(); err != nil {
		return nil, err
	}

	dom := &httputil.DOM{}
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.Header.Set("Cookie", c.cookie)
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for i, div := range dom.QueryAll("div", "class", "popular-container") {
		catalog := indextocatalog(i)
		if c.Option.Catalog == "" || c.Option.Catalog == catalog {
			for _, a := range div.QueryAll("a") {
				title := strings.TrimSpace(a.Text())
				url := "https://www.maoyan.com" + strings.TrimSpace(a.Href())
				if title == "" {
					continue
				}
				board.AppendTitleURLCatalog(title, url, catalog)
			}
		}
	}
	return board, nil
}

func (c *Crawler) updatecookie() error {
	return httputil.UpdateCookie("GET", "https://www.maoyan.com/", nil, httputil.NewOption(c.Option, ProxySwitch), &c.cookie)
}

func indextocatalog(i int) string {
	switch i {
	case 0:
		return Actor
	case 1:
		return News
	default:
		return Unknown
	}
}
