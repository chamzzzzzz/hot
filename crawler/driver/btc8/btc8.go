package btc8

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "btc8"
	ProxySwitch = false
	URL         = "https://www.8btc.com"
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
	div, err := dom.Find("div", "class", "hot-article-list")
	if err != nil {
		return nil, err
	}
	a, err := div.Find("a")
	if err != nil {
		return nil, err
	}
	p, err := a.Find("p")
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(p.Text())
	url := "https://www.8btc.com" + strings.TrimSpace(a.Href())
	board.AppendTitleURL(title, url)
	for _, a := range div.QueryAll("a", "class", "link-dark-major") {
		p, err = a.Find("p")
		if err != nil {
			return nil, err
		}
		title = strings.TrimSpace(p.Text())
		url = "https://www.8btc.com" + strings.TrimSpace(a.Href())
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
