package main

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/archiver/database"
	"github.com/chamzzzzzz/hot/crawler/baidu"
	"github.com/chamzzzzzz/hot/crawler/douyin"
	"github.com/chamzzzzzz/hot/crawler/github"
	"github.com/chamzzzzzz/hot/crawler/tieba"
	"github.com/chamzzzzzz/hot/crawler/toutiao"
	"github.com/chamzzzzzz/hot/crawler/v2ex"
	"github.com/chamzzzzzz/hot/crawler/weibo"
	"github.com/chamzzzzzz/hot/crawler/zhihu"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"time"
)

var (
	defaultLogger = log.New(os.Stdout, "collector: ", log.Ldate|log.Lmicroseconds)
	cronLogger    = log.New(os.Stdout, "cron: ", log.Ldate|log.Lmicroseconds)
)

type HotCollector struct {
	crawlers []hot.Crawler
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

	cookie := os.Getenv("HOT_COLLECT_WEIBO_COOKIE")
	if cookie == "" {
		return fmt.Errorf("missing env HOT_COLLECT_WEIBO_COOKIE")
	}

	proxy := os.Getenv("HOT_COLLECT_PROXY")
	if proxy == "" {
		return fmt.Errorf("missing env HOT_COLLECT_PROXY")
	}

	hc.archiver = &database.Archiver{
		DriverName:     driverName,
		DataSourceName: dataSourceName,
	}

	hc.crawlers = append(hc.crawlers, &baidu.Crawler{})
	hc.crawlers = append(hc.crawlers, &douyin.Crawler{})
	hc.crawlers = append(hc.crawlers, &toutiao.Crawler{})
	hc.crawlers = append(hc.crawlers, &weibo.Crawler{cookie})
	hc.crawlers = append(hc.crawlers, &zhihu.Crawler{})
	hc.crawlers = append(hc.crawlers, &v2ex.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &tieba.Crawler{})
	hc.crawlers = append(hc.crawlers, &github.Crawler{proxy})

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

	logger := cron.VerbosePrintfLogger(cronLogger)
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
	for _, crawler := range hc.crawlers {
		board, err := crawler.Crawl()
		if err != nil {
			defaultLogger.Printf("crawl, error='%s', crawler=%s\n", err, crawler.Name())
			continue
		}
		defaultLogger.Printf("crawl, crawler=%s, board=%s, count=%d\n", crawler.Name(), board.Name, len(board.Hots))

		var archived = 0
		if archived, err = hc.archiver.Archive(board); err != nil {
			defaultLogger.Printf("archive, error='%s', archiver=%s, board=%s, count=%d, archived=%d\n", err, hc.archiver.Name(), board.Name, len(board.Hots), archived)
			continue
		}
		defaultLogger.Printf("archive, archiver=%s, board=%s, count=%d, arvhiced=%d\n", hc.archiver.Name(), board.Name, len(board.Hots), archived)
	}
}

func main() {
	hc := &HotCollector{}
	err := hc.Start()
	if err != nil {
		defaultLogger.Printf("start, error='%s'\n", err)
		os.Exit(1)
	}
}
