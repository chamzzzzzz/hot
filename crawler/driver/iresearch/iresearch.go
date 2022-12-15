package iresearch

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "iresearch"
	ProxySwitch = false
	URL         = "https://www.iresearch.cn"
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
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.ContentEncoding = "gbk"
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div := dom.FindStrict("div", "class", "g-top f-mt-auto")
	if div.Error != nil {
		return nil, div.Error
	}
	div2 := div.Find("div", "class", "m-bd")
	if div2.Error != nil {
		return nil, div2.Error
	}
	for _, a := range div2.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
