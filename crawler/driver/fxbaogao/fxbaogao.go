package fxbaogao

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "fxbaogao"
	ProxySwitch = false
	URL         = "https://www.fxbaogao.com"
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
	for _, div := range dom.QueryAll("div", "class", "style_report__ZfacW") {
		a, err := div.Find("a")
		if err != nil {
			return nil, err
		}
		span, err := div.Find("span")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Text())
		url := "https://www.fxbaogao.com" + strings.TrimSpace(a.Href())
		date, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	for _, div := range dom.QueryAll("div", "class", "style_hotCardR__t9P0y") {
		//soup bug?
		//parsed duplicate a element in the div element.
		a, err := div.Find("a")
		if err != nil {
			return nil, err
		}
		p, err := div.Find("p")
		if err != nil {
			return nil, err
		}
		div, err = div.Find("div", "class", "style_time__bVTcg")
		if err != nil {
			return nil, err
		}
		span, err := div.Find("span")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(p.Text())
		url := "https://www.fxbaogao.com" + strings.TrimSpace(a.Href())
		date, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
