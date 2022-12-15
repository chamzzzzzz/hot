package timecom

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "timecom"
	ProxySwitch = true
	URL         = "https://time.com/"
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
	for _, li := range dom.FindAllStrict("li", "class", "most-popular-feed__item") {
		for i, a := range li.FindAll("a") {
			if i == 1 {
				h := a.Find("h3")
				if h.Error != nil {
					return nil, h.Error
				}
				title := strings.TrimSpace(h.Text())
				url := strings.TrimSpace(a.Attrs()["href"])
				board.AppendTitleURL(title, url)
			}
		}
	}
	return board, nil
}
