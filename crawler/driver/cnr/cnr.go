package cnr

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "cnr"
	ProxySwitch = false
	URL         = "http://news.cnr.cn/rebang/"
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
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.DetectContentEncoding = true
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, a := range dom.Query("div", "class", "articleList").QueryAll("a") {
		strong, err := a.Find("strong")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(strong.Text())
		url := strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{Title: title, URL: url})
	}
	return board, nil
}
