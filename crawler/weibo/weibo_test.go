package weibo

import (
	"os"
	"testing"
	"time"
)

func TestCrawl(t *testing.T) {
	c := Crawler{
		Cookie: os.Getenv("HOT_CRAWLER_WEIBO_TEST_COOKIE"),
	}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot.Title, hot.Date.Format(time.RFC3339))
		}
	}
}
