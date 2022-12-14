package jrj

import (
	"github.com/chamzzzzzz/hot/crawler/driver"
	"testing"
)

func TestCrawl(t *testing.T) {
	c := Crawler{}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlFinance(t *testing.T) {
	c := Crawler{Option: driver.Option{Catalog: Finance}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlTech(t *testing.T) {
	c := Crawler{Option: driver.Option{Catalog: Tech}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlHouse(t *testing.T) {
	c := Crawler{Option: driver.Option{Catalog: House}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
