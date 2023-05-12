package kanxue

import (
	"strings"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	News = "news"
	BBS  = "bbs"
)

const (
	DriverName  = "kanxue"
	ProxySwitch = false
	URL         = "https://bbs.pediy.com/thread-hotlist-all-0.htm"
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
	case News:
		return c.news()
	case BBS:
		return c.bbs()
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	b1, err := c.news()
	if err != nil {
		return nil, err
	}
	b2, err := c.bbs()
	if err != nil {
		return nil, err
	}
	for _, hot := range b1.Hots {
		board.Append(hot)
	}
	for _, hot := range b2.Hots {
		board.Append(hot)
	}
	return board, nil
}

func (c *Crawler) bbs() (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	tbody, err := dom.Find("tbody", "id", "arctilelist")
	if err != nil {
		return nil, err
	}
	for _, tr := range tbody.QueryAll("tr") {
		a, err := tr.Find("a", "class", "text-white")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Text())
		url := "https://bbs.pediy.com/" + strings.TrimSpace(a.Href())
		board.Append(&hot.Hot{
			Title:   title,
			URL:     url,
			Catalog: BBS,
		})
	}
	return board, nil
}

func (c *Crawler) news() (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", "https://www.kanxue.com", nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, div := range dom.QueryAll("div", "class", "pr-3 pb-2 mb-2 position-relative") {
		a, err := div.Find("a")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		if !strings.HasPrefix(url, "http") {
			url = "https://" + url
		}
		board.Append(&hot.Hot{
			Title:   title,
			URL:     url,
			Catalog: News,
		})
	}
	return board, nil
}
