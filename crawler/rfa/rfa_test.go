package rfa

import (
	"os"
	"testing"
)

func TestCrawlMandarin(t *testing.T) {
	c := Crawler{
		BoardName: Mandarin,
		Proxy:     os.Getenv("HOT_CRAWLER_TEST_PROXY"),
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
		BoardName: Cantonese,
		Proxy:     os.Getenv("HOT_CRAWLER_TEST_PROXY"),
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
		BoardName: English,
		Proxy:     os.Getenv("HOT_CRAWLER_TEST_PROXY"),
	}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
