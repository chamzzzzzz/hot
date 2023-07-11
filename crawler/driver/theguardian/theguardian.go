package theguardian

import (
	"fmt"
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "theguardian"
	ProxySwitch = true
	URL         = "https://www.theguardian.com/uk"
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
	for _, a := range dom.Query("ol", "class", "dcr-hqoq27").QueryAll("a") {
		h4, err := a.Find("h4", "class", "dcr-n4owjk")
		if err != nil {
			return nil, err
		}
		spans := h4.QueryAll("span")
		if len(spans) == 0 {
			return nil, fmt.Errorf("span not found")
		}
		span := spans[len(spans)-1]
		title := strings.TrimSpace(span.Text())
		url := strings.TrimSpace(a.Href())
		if !strings.HasPrefix(url, "http") {
			url = "https://www.theguardian.com" + url
		}
		board.Append(&hot.Hot{Title: title, URL: url})
	}
	return board, nil
}
