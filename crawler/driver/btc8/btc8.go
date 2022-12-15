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
	div := dom.Find("div", "class", "hot-article-list")
	if div.Error != nil {
		return nil, div.Error
	}
	a := div.Find("a")
	if a.Error != nil {
		return nil, a.Error
	}
	p := a.Find("p")
	if p.Error != nil {
		return nil, p.Error
	}
	title := strings.TrimSpace(p.Text())
	url := "https://www.8btc.com" + strings.TrimSpace(a.Attrs()["href"])
	board.AppendTitleURL(title, url)
	for _, a := range div.FindAllStrict("a", "class", "link-dark-major") {
		p = a.Find("p")
		if p.Error != nil {
			return nil, p.Error
		}
		title = strings.TrimSpace(p.Text())
		url = "https://www.8btc.com" + strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
