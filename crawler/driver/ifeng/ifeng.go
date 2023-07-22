package ifeng

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "ifeng"
	ProxySwitch = false
	URL         = "https://www.ifeng.com/"
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
	div, err := dom.Find("div", "class", "index_hot_box_Ztoic")
	if err != nil {
		return nil, err
	}
	h3, err := div.Find("h3", "class", "index_list_title_1x9s7 index_big_SGP4Y")
	if err != nil {
		return nil, err
	}
	a, err := h3.Find("a")
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(a.Text())
	url := strings.TrimSpace(a.Href())
	board.Append(&hot.Hot{Title: title, URL: url})
	for _, p := range div.QueryAll("p", "class", "index_news_list_p_5zOEF ") {
		a, err = p.Find("a")
		if err != nil {
			return nil, err
		}
		title = strings.TrimSpace(a.Text())
		url = strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{Title: title, URL: url})
	}
	return board, nil
}
