package main

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/archiver/database"
	"github.com/chamzzzzzz/hot/crawler"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"sync"
	"time"
)

var (
	logger = log.New(os.Stdout, "collector: ", log.Ldate|log.Lmicroseconds)
)

type HotCollector struct {
	option   crawler.Option
	crawlers []*crawler.Crawler
	archiver hot.Archiver
}

func (hc *HotCollector) Start() error {
	driverName := os.Getenv("HOT_COLLECT_DATABASE_DRIVER_NAME")
	dataSourceName := os.Getenv("HOT_COLLECTOR_DATABASE_DATA_SOURCE_NAME")
	if driverName == "" {
		return fmt.Errorf("missing env HOT_COLLECT_DATABASE_DRIVER_NAME")
	}
	if dataSourceName == "" {
		return fmt.Errorf("missing env HOT_COLLECTOR_DATABASE_DATA_SOURCE_NAME")
	}
	hc.archiver = &database.Archiver{
		DriverName:     driverName,
		DataSourceName: dataSourceName,
	}

	proxy := os.Getenv("HOT_COLLECT_PROXY")
	if proxy == "" {
		return fmt.Errorf("missing env HOT_COLLECT_PROXY")
	}
	hc.option.Proxy = proxy
	for _, driverName := range crawler.Drivers() {
		option := hc.option
		option.DriverName = driverName
		c, err := crawler.Open(option)
		if err != nil {
			return err
		}
		hc.crawlers = append(hc.crawlers, c)
	}

	spec := os.Getenv("HOT_COLLECT_CRON_SPEC")
	if spec == "" {
		spec = "5 * * * *"
	}

	tz := os.Getenv("HOT_COLLECT_CRON_TZ")
	if tz == "" {
		tz = "Local"
	}
	location, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}

	logger := cron.VerbosePrintfLogger(log.New(os.Stdout, "collector-cron: ", log.Ldate|log.Lmicroseconds))
	c := cron.New(
		cron.WithLocation(location),
		cron.WithLogger(logger),
		cron.WithChain(cron.SkipIfStillRunning(logger)),
	)

	c.AddJob(spec, hc)
	c.Run()
	return nil
}

func (hc *HotCollector) Run() {
	t1 := time.Now()
	var wg sync.WaitGroup
	for _, crawler := range hc.crawlers {
		wg.Add(1)
		go func(crawler hot.Crawler) {
			defer wg.Done()
			board, err := crawler.Crawl()
			if err != nil {
				logger.Printf("crawl, error='%s', crawler=%s\n", err, crawler.Name())
				return
			}
			logger.Printf("crawl, crawler=%s, board=%s, count=%d\n", crawler.Name(), board.Name, len(board.Hots))

			var archived = 0
			if archived, err = hc.archiver.Archive(board); err != nil {
				logger.Printf("archive, error='%s', archiver=%s, board=%s, count=%d, archived=%d\n", err, hc.archiver.Name(), board.Name, len(board.Hots), archived)
				return
			}
			logger.Printf("archive, archiver=%s, board=%s, count=%d, arvhiced=%d\n", hc.archiver.Name(), board.Name, len(board.Hots), archived)
		}(crawler)
	}
	wg.Wait()
	logger.Printf("run, used=%v", time.Since(t1))
}

func main() {
	hc := &HotCollector{}
	err := hc.Start()
	if err != nil {
		logger.Printf("start, error='%s'\n", err)
		os.Exit(1)
	}
}
