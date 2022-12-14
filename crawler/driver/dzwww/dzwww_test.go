package dzwww

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

func TestCrawlHotSearch(t *testing.T) {
	c := Crawler{Option: driver.Option{Catalog: HotSearch}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlHotNews(t *testing.T) {
	c := Crawler{Option: driver.Option{Catalog: HotNews}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
