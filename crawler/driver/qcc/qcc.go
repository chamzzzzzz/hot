package qcc

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "qcc"
	ProxySwitch = false
	URL         = "https://www.qcc.com/top_search"
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
	for _, div := range dom.FindAll("div", "class", "hslist-item") {
		a := div.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}

func (c *Crawler) updatecookie() error {
	return httputil.UpdateCookie("GET", "https://www.qcc.com", nil, httputil.NewOption(c.Option, ProxySwitch), &c.cookie)
}
