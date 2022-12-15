package credit51

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "credit51"
	ProxySwitch = false
	URL         = "https://bbs.51credit.com"
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
	div := dom.Find("div", "class", "hot-post")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a") {
		span := a.FindStrict("span", "class", "right text-deal")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(span.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
