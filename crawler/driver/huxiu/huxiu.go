package huxiu

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "huxiu"
	ProxySwitch = false
	URL         = "https://www.huxiu.com"
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
	for _, a := range dom.FindAllStrict("a", "class", "lirt-item__right__content") {
		h3 := a.Find("h3")
		if h3.Error != nil {
			return nil, h3.Error
		}
		title := strings.TrimSpace(h3.Text())
		url := "https://www.huxiu.com" + strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}