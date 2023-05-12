package mddcloud

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "mddcloud"
	ProxySwitch = false
	URL         = "https://www.mddcloud.com.cn/"
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
	for _, li := range dom.QueryAll("li", "class", "rank-list-item") {
		a, err := li.Find("a", "class", "g-a-block")
		if err != nil {
			return nil, err
		}
		p1, err := li.Find("p", "class", "rank-vod-title")
		if err != nil {
			return nil, err
		}
		p2, err := li.Find("p", "class", "rank-vod-desc")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(p1.Text())
		summary := strings.TrimSpace(p2.Text())
		url := "https://www.mddcloud.com.cn" + strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{Title: title, Summary: summary, URL: url})
	}
	return board, nil
}
