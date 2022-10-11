package bbc

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
