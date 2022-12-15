package github

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
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
	for _, article := range dom.FindAllStrict("article", "class", "Box-row") {
		h1 := article.FindStrict("h1", "class", "h3 lh-condensed")
		if h1.Error != nil {
			continue
		}

		a := h1.FindStrict("a")
		if a.Error != nil {
			continue
		}

		p := article.FindStrict("p", "class", "col-9 color-fg-muted my-1 pr-4")
		if p.Error != nil {
			continue
		}

		title := strings.Trim(a.Attrs()["href"], "/")
		summary := strings.Trim(p.Text(), " \n")
		url := "https://github.com/" + title
		board.AppendTitleSummaryURL(title, summary, url)
	}
	return board, nil
}
