package jrj

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	Finance = "finance"
	Tech    = "tech"
	House   = "house"
)

const (
	DriverName  = "jrj"
	ProxySwitch = false
	URL         = "http://tech.jrj.com.cn"
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
	switch c.Option.Catalog {
	case Finance:
		return c.finance()
	case Tech:
		return c.tech()
	case House:
		return c.house()
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	if b, err := c.finance(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}

	if b, err := c.tech(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}

	if b, err := c.house(); err != nil {
		return nil, err
	} else {
		for _, hot := range b.Hots {
			board.Append0(hot)
		}
	}
	return board, nil
}

func (c *Crawler) tech() (*hot.Board, error) {
	dom := &httputil.DOM{}
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.ContentEncoding = "gbk"
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div, err := dom.Find("div", "class", "hotart hotnews")
	if err != nil {
		return nil, err
	}
	for _, a := range div.QueryAll("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		board.Append4(title, "", url, Tech)
	}
	return board, nil
}

func (c *Crawler) house() (*hot.Board, error) {
	dom := &httputil.DOM{}
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.ContentEncoding = "gbk"
	if err := httputil.Request("GET", "http://house.jrj.com.cn", nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div, err := dom.Find("div", "class", "hotart")
	if err != nil {
		return nil, err
	}
	for _, a := range div.QueryAll("a") {
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		board.Append4(title, "", url, House)
	}
	return board, nil
}

func (c *Crawler) finance() (*hot.Board, error) {
	dom := &httputil.DOM{}
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.ContentEncoding = "gbk"
	if err := httputil.Request("GET", "http://finance.jrj.com.cn/list/industrynews.shtml", nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, id := range []string{"con_1", "con_2"} {
		ul, err := dom.Find("ul", "id", id)
		if err != nil {
			return nil, err
		}
		for _, a := range ul.QueryAll("a") {
			title := strings.TrimSpace(a.Text())
			url := strings.TrimSpace(a.Href())
			board.Append4(title, "", url, Finance)
		}
	}
	return board, nil
}
