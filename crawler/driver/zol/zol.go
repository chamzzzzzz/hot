package zol

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "zol"
	ProxySwitch = false
	URL         = "https://news.zol.com.cn/"
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
	for _, id := range []string{"list-v-1", "list-v-2"} {
		div, err := dom.Find("div", "id", id)
		if err != nil {
			return nil, err
		}
		for _, a := range div.QueryAll("a") {
			title := strings.TrimSpace(a.Text())
			url := "https:" + strings.TrimSpace(a.Href())
			board.Append(&hot.Hot{Title: title, URL: url})
		}
	}
	return board, nil
}
