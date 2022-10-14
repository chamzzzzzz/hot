package rfa

import (
	"os"
	"testing"
)

func TestCrawl(t *testing.T) {
	c := Crawler{
		Proxy: os.Getenv("HOT_CRAWLER_TEST_PROXY"),
	}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlMandarin(t *testing.T) {
	c := Crawler{
		Catalog: Mandarin,
		Proxy:   os.Getenv("HOT_CRAWLER_TEST_PROXY"),
	}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlCantonese(t *testing.T) {
	c := Crawler{
		Catalog: Cantonese,
		Proxy:   os.Getenv("HOT_CRAWLER_TEST_PROXY"),
	}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}

func TestCrawlEnglish(t *testing.T) {
	c := Crawler{
		Catalog: English,
		Proxy:   os.Getenv("HOT_CRAWLER_TEST_PROXY"),
	}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
