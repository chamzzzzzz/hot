package weibo

import (
	"github.com/chamzzzzzz/hot/crawler/driver"
	"testing"
)

func TestCrawl(t *testing.T) {
	c := Crawler{Option: driver.Option{}}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
