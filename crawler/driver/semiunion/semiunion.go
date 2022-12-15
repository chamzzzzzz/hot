package semiunion

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "semiunion"
	ProxySwitch = false
	URL         = "http://www.semiunion.com/insight/"
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
	for _, li := range dom.FindAllStrict("li", "class", "each-news") {
		div := li.Find("div", "class", "name")
		if div.Error != nil {
			return nil, div.Error
		}
		a := div.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		div2 := li.Find("div", "class", "desc")
		if div2.Error != nil {
			return nil, div2.Error
		}
		span := li.Find("span", "class", "time")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(a.Text())
		summary := strings.TrimSpace(div2.Text())
		url := "http://www.semiunion.com" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(span.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}
