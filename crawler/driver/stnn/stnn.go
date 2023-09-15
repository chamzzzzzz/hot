package stnn

import (
	"fmt"
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "stnn"
	ProxySwitch = false
	URL         = "https://www.stnn.cc"
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
	divs := dom.QueryAll("div", "class", "rightitem")
	if len(divs) != 4 {
		return nil, fmt.Errorf("divs length is not 4")
	}
	for _, a := range divs[2].QueryAll("a") {
		title := strings.TrimSpace(a.Text())
		url := "https://www.stnn.cc" + strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{Title: title, URL: url})
	}
	return board, nil
}
