package nytimes

import (
	"fmt"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "nytimes"
	ProxySwitch = true
	URL         = "https://cn.nytimes.com/async/mostviewed/all/?lang=zh-hans"
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
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}
	if body.Code != 0 {
		return nil, fmt.Errorf("body code: %d", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, daily := range body.List.Daily {
		board.Append(&hot.Hot{Title: daily.Headline, Summary: daily.Summary, URL: daily.URL})
	}
	return board, nil
}

type body struct {
	Code int `json:"code"`
	List struct {
		Daily []struct {
			URL           string `json:"url"`
			Headline      string `json:"headline"`
			ShortHeadline string `json:"short_headline"`
			Summary       string `json:"summary"`
		} `json:"daily"`
		Weekly []struct {
			URL           string `json:"url"`
			Headline      string `json:"headline"`
			ShortHeadline string `json:"short_headline"`
			Summary       string `json:"summary"`
		} `json:"weekly"`
		WeeklySlideshow []struct {
			URL           string `json:"url"`
			Headline      string `json:"headline"`
			ShortHeadline string `json:"short_headline"`
			Summary       string `json:"summary"`
		} `json:"weekly_slideshow"`
	} `json:"list"`
}
