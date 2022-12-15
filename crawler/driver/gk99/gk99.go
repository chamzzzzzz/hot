package gk99

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "gk99"
	ProxySwitch = false
	URL         = "http://dota2.gk99.com/rd/"
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
	ul := dom.FindStrict("ul", "class", "cfix fin_newsList")
	if ul.Error != nil {
		return nil, ul.Error
	}
	for _, li := range ul.FindAllStrict("li") {
		a := li.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		p := li.Find("p")
		if p.Error != nil {
			return nil, p.Error
		}
		em := li.Find("em", "class", "fRight")
		if em.Error != nil {
			return nil, em.Error
		}
		title := strings.TrimSpace(a.Text())
		summary := strings.TrimSpace(p.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02 15:04:05", strings.TrimSpace(em.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}