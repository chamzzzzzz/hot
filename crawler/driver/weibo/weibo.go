package weibo

import (
	"strconv"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "weibo"
	ProxySwitch = false
	URL         = "https://s.weibo.com/top/summary"
	Cookie      = "WBPSESS=durPiJxsbzq5XDaI2wW0NxQldYOrBwQzLVlPfvpcy3mQ3XQonV49sfubFlqvuI_rBrarQ2dZHLfrOVaRKnvrm9130Jsv26CGHwu2LjHl3RrnHDHKIUtZPYEi9qKk6n-K; SUB=_2AkMU1LJTf8NxqwJRmPAQymrhaYl_yg_EieKiiEOIJRMxHRl-yT92qkI6tRB6P1ScvMDt8JtdZqvVJlNftBcRg-WjvODv; SUBP=0033WrSXqPxfM72-Ws9jqgMF55529P9D9WFSWJI0b_93sKJGpCc_.aOL; XSRF-TOKEN=Z3qrKi3V9M0TVao6eTMMmpRC"
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
	cookie string
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	if err := c.updatecookie(); err != nil {
		return nil, err
	}

	dom := &httputil.DOM{}
	option := httputil.NewOption(c.Option, ProxySwitch)
	option.Header.Set("Cookie", c.cookie)
	if err := httputil.Request("GET", URL, nil, "dom", dom, option); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div, err := dom.Find("div", "id", "pl_top_realtimehot")
	if err != nil {
		return nil, err
	}
	tbody, err := div.Find("tbody")
	if err != nil {
		return nil, err
	}
	for _, tr := range tbody.QueryAll("tr", "class", "") {
		td01, err := tr.Find("td", "class", "td-01")
		if err != nil {
			return nil, err
		}
		if _, err := strconv.Atoi(td01.Text()); err != nil {
			continue
		}
		td02, err := tr.Find("td", "class", "td-02")
		if err != nil {
			return nil, err
		}
		a, err := td02.Find("a")
		if err != nil {
			return nil, err
		}
		board.Append(&hot.Hot{
			Title: a.Text(),
		})
	}
	return board, nil
}

func (c *Crawler) updatecookie() error {
	if c.cookie != "" {
		return nil
	}
	if c.Option.Cookie != "" {
		c.cookie = c.Option.Cookie
		return nil
	}
	c.cookie = Cookie
	return nil
}
