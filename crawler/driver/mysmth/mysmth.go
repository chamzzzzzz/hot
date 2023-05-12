package mysmth

import (
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "mysmth"
	ProxySwitch = false
	URL         = "https://www.mysmth.net/nForum/rss/topten"
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
	if err := httputil.Request("GET", URL, nil, "xml", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, item := range body.Channel.Item {
		title := item.Title
		url := item.Link
		date, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return nil, err
		}
		if title == "" {
			continue
		}
		board.Append(&hot.Hot{Title: title, URL: url, PublishDate: date})
	}
	return board, nil
}

type body struct {
	Channel struct {
		Item []struct {
			Title   string `xml:"title"`
			Link    string `xml:"link"`
			PubDate string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (body *body) NormalizedCode() int {
	return 0
}
