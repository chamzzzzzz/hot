package ebrun

import (
	"testing"

	"github.com/chamzzzzzz/hot/crawler/driver"
)

func TestCrawl(t *testing.T) {
	c := Crawler{Option: driver.NewTestOptionFromEnv()}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot)
		}
	}
}
