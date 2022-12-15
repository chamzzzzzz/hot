package sputniknews

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strconv"
	"strings"
	"time"
)

const (
	DriverName  = "sputniknews"
	ProxySwitch = false
	URL         = "https://www.sputniknews.cn"
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
	div := dom.Find("div", "data-section", "3")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a", "class", "cell-list__item m-no-image") {
		span := a.Find("span", "class", "cell__controls-date")
		if span.Error != nil {
			return nil, span.Error
		}
		span = span.Find("span")
		if span.Error != nil {
			return nil, span.Error
		}

		timestamp, err := strconv.ParseInt(span.Attrs()["data-unixtime"], 10, 64)
		if err != nil {
			return nil, err
		}

		title := strings.TrimSpace(a.Attrs()["title"])
		url := "https://www.sputniknews.cn" + strings.TrimSpace(a.Attrs()["href"])
		date := time.Unix(timestamp, 0)
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
