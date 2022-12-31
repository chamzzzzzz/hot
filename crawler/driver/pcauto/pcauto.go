package pcauto

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	Article = "article"
	Topic   = "topic"
	Forum   = "forum"
	Unknown = "unknown"
)

const (
	DriverName  = "pcauto"
	ProxySwitch = false
	URL         = "https://www.pcauto.com.cn/"
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
	case Article:
		return c.article()
	case Topic:
		return c.topic()
	case Forum:
		return c.forum()
	default:
		return c.all()
	}
}

func (c *Crawler) article() (*hot.Board, error) {
	return c.withClasses("txts")
}

func (c *Crawler) topic() (*hot.Board, error) {
	return c.withClasses("bbs_topics")
}

func (c *Crawler) forum() (*hot.Board, error) {
	return c.withClasses("hotForums")
}

func (c *Crawler) all() (*hot.Board, error) {
	return c.withClasses("txts", "bbs_topics", "hotForums")
}

func (c *Crawler) withClasses(classes ...string) (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", URL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	div, err := dom.Find("div", "class", "section ranking")
	if err != nil {
		return nil, err
	}
	for _, class := range classes {
		ul, err := div.Find("ul", "class", class)
		if err != nil {
			return nil, err
		}
		catalog := classtocatalog(class)
		for _, a := range ul.QueryAll("a") {
			title := strings.TrimSpace(a.Text())
			url := "https:" + strings.TrimSpace(a.Href())
			board.AppendTitleURLCatalog(title, url, catalog)
		}
	}
	return board, nil
}

func classtocatalog(class string) string {
	switch class {
	case "txts":
		return Article
	case "bbs_topics":
		return Topic
	case "hotForums":
		return Forum
	default:
		return Unknown
	}
}
