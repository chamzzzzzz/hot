package wsj

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "wsj"
	ProxySwitch = true
	URL         = "https://cn.wsj.com/zh-hans/rss"
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
	body := &body{}
	if err := httputil.Request("GET", URL, nil, "xml", &body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, item := range body.Channel.Item {
		title := item.Title
		summary := item.Description
		url := item.Link
		date, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return nil, err
		}
		date = time.Unix(date.Unix(), 0)
		i := strings.Index(summary, "<p>")
		j := strings.Index(summary, "</p>")
		if i >= 0 && j >= 0 {
			summary = summary[i+3 : j]
		}
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	Text    string `xml:",chardata"`
	Channel struct {
		Item []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}
