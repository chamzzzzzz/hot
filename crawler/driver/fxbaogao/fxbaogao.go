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
	for _, div := range dom.FindAllStrict("div", "class", "style_report__ZfacW") {
		a := div.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		span := div.Find("span")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(a.Text())
		url := "https://www.fxbaogao.com" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	for _, div := range dom.FindAllStrict("div", "class", "style_hotCardR__t9P0y") {
		//soup bug?
		//parsed duplicate a element in the div element.
		a := div.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		p := div.Find("p")
		if p.Error != nil {
			return nil, p.Error
		}
		div = div.Find("div", "class", "style_time__bVTcg")
		if div.Error != nil {
			return nil, div.Error
		}
		span := div.Find("span")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(p.Text())
		url := "https://www.fxbaogao.com" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}