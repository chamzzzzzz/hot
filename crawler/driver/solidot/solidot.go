package solidot

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "solidot"
	ProxySwitch = false
	URL         = "https://www.solidot.org"
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
	ul := dom.Find("ul", "class", "old_articles")
	if ul.Error != nil {
		return nil, ul.Error
	}
	ul2 := ul.Find("ul", "class", "comment_on")
	if ul2.Error != nil {
		return nil, ul2.Error
	}
	for _, a := range ul2.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		url := "https://www.solidot.org" + strings.TrimSpace(a.Attrs()["href"])
		board.AppendTitleURL(title, url)
	}
	return board, nil
}
