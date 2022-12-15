package bbc

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "bbc"
	ProxySwitch = true
	URL         = "https://www.bbc.com/zhongwen/mostread/simp.json"
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

	board := hot.NewBoard(c.Name())
	for _, record := range body.Records {
		title := strings.TrimSpace(record.Promo.Headlines.ShortHeadline)
		summary := strings.TrimSpace(record.Promo.Summary)
		url := "https://www.bbc.com/" + strings.TrimSpace(strings.Trim(record.Promo.ID, "urn:bbc:ares::asset:"))
		date := time.UnixMilli(record.Promo.Timestamp)
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	Records []struct {
		Promo struct {
			Headlines struct {
				ShortHeadline string `json:"shortHeadline"`
				Headline      string `json:"headline"`
			} `json:"headlines"`
			Summary   string `json:"summary"`
			Timestamp int64  `json:"timestamp"`
			ID        string `json:"id"`
		} `json:"promo,omitempty"`
	} `json:"records"`
}
