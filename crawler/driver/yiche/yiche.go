package yiche

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	DriverName  = "yiche"
	ProxySwitch = false
	URL         = "https://www.yiche.com"
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
	for _, class := range []string{"wenzhang_l ka", "retie_l ka", "hotshipin_l"} {
		for _, li := range dom.FindAllStrict("li", "class", class) {
			a := li.Find("a")
			if a.Error != nil {
				return nil, a.Error
			}
			title := strings.TrimSpace(a.Text())
			url := "https:" + strings.TrimSpace(a.Attrs()["href"])
			board.AppendTitleURL(title, url)
		}
	}
	return board, nil
}
