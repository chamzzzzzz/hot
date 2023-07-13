package takungpao

import (
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "takungpao"
	ProxySwitch = false
	URL         = "http://www.takungpao.com"
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
	div, err := dom.Find("div", "class", "ranking")
	if err != nil {
		return nil, err
	}
	for _, dl := range div.QueryAll("dl") {
		a, err := dl.Find("a")
		if err != nil {
			return nil, err
		}
		dd, err := dl.Find("dd", "class", "time")
		if err != nil {
			return nil, err
		}

		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		date, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(dd.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append(&hot.Hot{Title: title, URL: url, PublishDate: date})
	}
	return board, nil
}
