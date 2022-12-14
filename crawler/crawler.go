package crawler

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
)

type (
	Option = driver.Option
)

type Crawler struct {
	dc     driver.Crawler
	option Option
}

func (c *Crawler) Driver() driver.Driver {
	return c.dc.Driver()
}

func (c *Crawler) Option() Option {
	return c.option
}

func (c *Crawler) Name() string {
	return c.dc.Name()
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	return c.dc.Crawl()
}

func Open(option Option) (*Crawler, error) {
	c, err := driver.Open(option)
	if err != nil {
		return nil, err
	}
	return &Crawler{c, option}, nil
}

func Drivers() []string {
	return driver.Drivers()
}
