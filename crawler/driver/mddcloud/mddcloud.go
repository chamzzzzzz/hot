package mddcloud

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
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
	for _, li := range dom.FindAllStrict("li", "class", "rank-list-item") {
		a := li.FindStrict("a", "class", "g-a-block")
		p1 := li.FindStrict("p", "class", "rank-vod-title")
		p2 := li.FindStrict("p", "class", "rank-vod-desc")
		if a.Error != nil {
			return nil, a.Error
		}
		if p1.Error != nil {
			return nil, p1.Error
		}
		if p2.Error != nil {
			return nil, p2.Error
		}
		title := strings.TrimSpace(p1.Text())
		summary := strings.TrimSpace(p2.Text())
		url := "https://www.mddcloud.com.cn" + strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleSummaryURL(title, summary, url)
	}
	return board, nil
}
