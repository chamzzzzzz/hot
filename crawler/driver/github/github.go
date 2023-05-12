package github

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "github"
	ProxySwitch = false
	URL         = "https://github.com/trending"
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
	for _, article := range dom.QueryAll("article", "class", "Box-row") {
		h2, err := article.Find("h2", "class", "h3 lh-condensed")
		if err != nil {
			continue
		}
		a, err := h2.Find("a")
		if err != nil {
			continue
		}
		p, err := article.Find("p", "class", "col-9 color-fg-muted my-1 pr-4")
		if err != nil {
			continue
		}
		title := strings.Trim(a.Href(), "/")
		summary := strings.Trim(p.Text(), " \n")
		url := "https://github.com/" + title
		board.Append(&hot.Hot{Title: title, Summary: summary, URL: url})
	}
	return board, nil
}
