package cnblogs

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "cnblogs"
	ProxySwitch = false
	URL         = "https://www.cnblogs.com/aggsite/SideRight"
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
	for _, ul := range dom.FindAllStrict("ul", "class", "item-list") {
		for _, a := range ul.FindAll("a") {
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Attrs()["href"])
			board.AppendTitleURL(title, url)
		}
	}
	return board, nil
}
