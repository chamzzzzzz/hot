package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/chamzzzzzz/hot/archiver/file"
	"github.com/chamzzzzzz/hot/crawler"
)

var (
	proxy    = os.Getenv("HOT_ARCHIVER_PROXY")
	archiver = &file.Archiver{}
	crawlers []*crawler.Crawler
)

func main() {
	log.Printf("proxy=%s\n", proxy)
	log.Printf("archiver=%s\n", archiver.Name())
	for _, driverName := range crawler.Drivers() {
		c, err := crawler.Open(crawler.Option{DriverName: driverName, Proxy: proxy})
		if err != nil {
			log.Printf("[%s] open crawler failed, err=%s\n", driverName, err)
			return
		}
		crawlers = append(crawlers, c)
	}
	for {
		archive()
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		log.Printf("next archive at %s\n", next.Format("2006-01-02 15:04:05"))
		time.Sleep(next.Sub(now))
	}
}

func archive() {
	log.Printf("start archive at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	t := time.Now()
	var wg sync.WaitGroup
	for _, c := range crawlers {
		wg.Add(1)
		go func(c *crawler.Crawler) {
			defer wg.Done()
			board, err := c.Crawl()
			if err != nil {
				log.Printf("[%s] crawl failed, err=%s\n", c.Name(), err)
				return
			}
			archiver.Archive(board)
		}(c)
	}
	wg.Wait()
	log.Printf("archive used %v\n", time.Since(t))
	log.Printf("finish archive at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
