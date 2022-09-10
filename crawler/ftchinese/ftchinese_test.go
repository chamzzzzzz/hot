package ftchinese

import (
	"os"
	"testing"
	"time"
)

func TestCrawl(t *testing.T) {
	c := Crawler{
		Proxy: os.Getenv("HOT_CRAWLER_TEST_PROXY"),
	}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot.Title, hot.Summary, hot.Date.Format(time.RFC3339))
		}
	}
}
