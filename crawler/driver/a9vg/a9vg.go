package a9vg

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "a9vg"
	ProxySwitch = false
	URL         = "https://www.a9vg.com/list/strategy"
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
	for _, ul := range dom.FindAllStrict("ul", "class", "vd-list a9-mini-news-list vdp-theme_lite is-gap is-divided") {
		for _, a := range ul.FindAllStrict("a", "class", "vd-list_container") {
			span := a.Find("span", "class", "a9-mini-news-list_title")
			if span.Error != nil {
				return nil, span.Error
			}
			title := strings.TrimSpace(span.Text())
			url := "https://www.a9vg.com" + strings.TrimSpace(a.Attrs()["href"])
			board.AppendTitleURL(title, url)
		}
	}
	return board, nil
}
