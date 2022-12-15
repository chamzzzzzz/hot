package im2maker

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "im2maker"
	ProxySwitch = false
	URL         = "https://www.im2maker.com"
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
	div := dom.Find("div", "id", "hot_posts_position")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, div2 := range div.FindAllStrict("div", "class", "desc") {
		a := div2.Find("a", "class", "title")
		if a.Error != nil {
			return nil, a.Error
		}
		span := div2.Find("span", "class", "timeago")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(span.Attrs()["datetime"]), time.Local)
		if err != nil {
			return nil, err
		}
		board.AppendTitleURLDate(title, url, date)
	}
	return board, nil
}
