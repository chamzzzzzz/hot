package rfa

import (
	"github.com/chamzzzzzz/hot/crawler/driver"
	"os"
	"testing"
)

func TestCrawl(t *testing.T) {
	c := Crawler{Option: driver.Option{Proxy: os.Getenv("HOT_CRAWLER_TEST_PROXY")}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlMandarin(t *testing.T) {
	c := Crawler{Option: driver.Option{Proxy: os.Getenv("HOT_CRAWLER_TEST_PROXY"), Catalog: Mandarin}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlCantonese(t *testing.T) {
	c := Crawler{Option: driver.Option{Proxy: os.Getenv("HOT_CRAWLER_TEST_PROXY"), Catalog: Cantonese}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlEnglish(t *testing.T) {
	c := Crawler{Option: driver.Option{Proxy: os.Getenv("HOT_CRAWLER_TEST_PROXY"), Catalog: English}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
