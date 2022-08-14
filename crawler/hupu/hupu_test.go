package hupu

import (
	"testing"
	"time"
)

func TestCrawl(t *testing.T) {
	c := Crawler{}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot.Title, hot.Summary, hot.Date.Format(time.RFC3339))
		}
	}
}

func TestCrawlBasketball(t *testing.T) {
	c := Crawler{Basketball}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot.Title, hot.Summary, hot.Date.Format(time.RFC3339))
		}
	}
}

func TestCrawlFootball(t *testing.T) {
	c := Crawler{Football}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot.Title, hot.Summary, hot.Date.Format(time.RFC3339))
		}
	}
}

func TestCrawlGambia(t *testing.T) {
	c := Crawler{Gambia}
	if board, err := c.Crawl(); err != nil {
		t.Error(err)
	} else {
		for _, hot := range board.Hots {
			t.Log(hot.Title, hot.Summary, hot.Date.Format(time.RFC3339))
		}
	}
}
