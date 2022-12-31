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
	div, err := dom.Find("div", "id", "hot_posts_position")
	if err != nil {
		return nil, err
	}
	for _, div2 := range div.QueryAll("div", "class", "desc") {
		a, err := div2.Find("a", "class", "title")
		if err != nil {
			return nil, err
		}
		span, err := div2.Find("span", "class", "timeago")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		date, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(span.Attribute("datetime")), time.Local)
		if err != nil {
			return nil, err
		}
		board.AppendTitleURLDate(title, url, date)
	}
	return board, nil
}
