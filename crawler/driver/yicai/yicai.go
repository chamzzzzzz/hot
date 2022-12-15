package yicai

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "yicai"
	ProxySwitch = false
	URL         = "https://www.yicai.com"
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
	div := dom.Find("div", "id", "hotest")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a") {
		h2 := a.Find("h2")
		if h2.Error != nil {
			return nil, h2.Error
		}
		title := strings.TrimSpace(h2.Text())
		url := "https://www.yicai.com" + strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
