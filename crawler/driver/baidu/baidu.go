package baidu

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "baidu"
	ProxySwitch = false
	URL         = "https://top.baidu.com/board?tab=realtime"
)

type Driver struct{}

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
	for _, div := range dom.QueryAll("div", "class", "content_1YWBm") {
		div01, err := div.Find("div", "class", "c-single-text-ellipsis")
		if err != nil {
			return nil, err
		}
		div02, err := div.Find("div", "class", "small_Uvkd3")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(div01.Text())
		summary := strings.TrimSpace(div02.Text())
		board.Append(&hot.Hot{Title: title, Summary: summary})
	}
	return board, nil
}
