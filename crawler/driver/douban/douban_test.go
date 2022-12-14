package douban

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

func TestCrawlNote(t *testing.T) {
	c := Crawler{Option: driver.Option{Catalog: Note}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlMovie(t *testing.T) {
	c := Crawler{Option: driver.Option{Catalog: Movie}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
