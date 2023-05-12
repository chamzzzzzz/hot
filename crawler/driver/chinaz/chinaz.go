package chinaz

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "chinaz"
	ProxySwitch = false
	URL         = "https://www.chinaz.com/"
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
	div, err := dom.Find("div", "class", "hots-wrapper__bd")
	if err != nil {
		return nil, err
	}
	for _, a := range div.QueryAll("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{Title: title, URL: url})
	}
	return board, nil
}

func (c *Crawler) updatecookie() error {
	return httputil.UpdateCookie("GET", URL, nil, httputil.NewOption(c.Option, ProxySwitch), &c.cookie)
}
