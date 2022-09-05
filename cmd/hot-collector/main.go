package main

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/archiver/database"
	"github.com/chamzzzzzz/hot/crawler/baidu"
	"github.com/chamzzzzzz/hot/crawler/bbc"
	"github.com/chamzzzzzz/hot/crawler/bilibili"
	"github.com/chamzzzzzz/hot/crawler/bjnews"
	"github.com/chamzzzzzz/hot/crawler/btc8"
	"github.com/chamzzzzzz/hot/crawler/chinanews"
	"github.com/chamzzzzzz/hot/crawler/chinaz"
	"github.com/chamzzzzzz/hot/crawler/chiphell"
	"github.com/chamzzzzzz/hot/crawler/cls"
	"github.com/chamzzzzzz/hot/crawler/cnbeta"
	"github.com/chamzzzzzz/hot/crawler/cnblogs"
	"github.com/chamzzzzzz/hot/crawler/cninfo"
	"github.com/chamzzzzzz/hot/crawler/csdn"
	"github.com/chamzzzzzz/hot/crawler/cto51"
	"github.com/chamzzzzzz/hot/crawler/daniu"
	"github.com/chamzzzzzz/hot/crawler/donews"
	"github.com/chamzzzzzz/hot/crawler/dongchedi"
	"github.com/chamzzzzzz/hot/crawler/douban"
	"github.com/chamzzzzzz/hot/crawler/douyin"
	"github.com/chamzzzzzz/hot/crawler/eastmoney"
	"github.com/chamzzzzzz/hot/crawler/fxbaogao"
	"github.com/chamzzzzzz/hot/crawler/gameres"
	"github.com/chamzzzzzz/hot/crawler/gamersky"
	"github.com/chamzzzzzz/hot/crawler/gelonghui"
	"github.com/chamzzzzzz/hot/crawler/gitchat"
	"github.com/chamzzzzzz/hot/crawler/github"
	"github.com/chamzzzzzz/hot/crawler/globaltimes"
	"github.com/chamzzzzzz/hot/crawler/haokan"
	"github.com/chamzzzzzz/hot/crawler/hibor"
	"github.com/chamzzzzzz/hot/crawler/hupu"
	"github.com/chamzzzzzz/hot/crawler/ifeng"
	"github.com/chamzzzzzz/hot/crawler/infoq"
	"github.com/chamzzzzzz/hot/crawler/investing"
	"github.com/chamzzzzzz/hot/crawler/ithome"
	"github.com/chamzzzzzz/hot/crawler/jin10"
	"github.com/chamzzzzzz/hot/crawler/jinse"
	"github.com/chamzzzzzz/hot/crawler/jisilu"
	"github.com/chamzzzzzz/hot/crawler/jqka10"
	"github.com/chamzzzzzz/hot/crawler/kanxue"
	"github.com/chamzzzzzz/hot/crawler/kr36"
	"github.com/chamzzzzzz/hot/crawler/kuaishou"
	"github.com/chamzzzzzz/hot/crawler/leikeji"
	"github.com/chamzzzzzz/hot/crawler/mydrivers"
	"github.com/chamzzzzzz/hot/crawler/netease"
	"github.com/chamzzzzzz/hot/crawler/nowcoder"
	"github.com/chamzzzzzz/hot/crawler/nytimes"
	"github.com/chamzzzzzz/hot/crawler/odaily"
	"github.com/chamzzzzzz/hot/crawler/oschina"
	"github.com/chamzzzzzz/hot/crawler/panews"
	"github.com/chamzzzzzz/hot/crawler/pearvideo"
	"github.com/chamzzzzzz/hot/crawler/pojie52"
	"github.com/chamzzzzzz/hot/crawler/readhub"
	"github.com/chamzzzzzz/hot/crawler/rfa"
	"github.com/chamzzzzzz/hot/crawler/so360"
	"github.com/chamzzzzzz/hot/crawler/sogou"
	"github.com/chamzzzzzz/hot/crawler/sohu"
	"github.com/chamzzzzzz/hot/crawler/solidot"
	"github.com/chamzzzzzz/hot/crawler/sspai"
	"github.com/chamzzzzzz/hot/crawler/techweb"
	"github.com/chamzzzzzz/hot/crawler/thecover"
	"github.com/chamzzzzzz/hot/crawler/thepaper"
	"github.com/chamzzzzzz/hot/crawler/tianya"
	"github.com/chamzzzzzz/hot/crawler/tieba"
	"github.com/chamzzzzzz/hot/crawler/timecom"
	"github.com/chamzzzzzz/hot/crawler/toutiao"
	"github.com/chamzzzzzz/hot/crawler/toutiaoio"
	"github.com/chamzzzzzz/hot/crawler/v2ex"
	"github.com/chamzzzzzz/hot/crawler/weibo"
	"github.com/chamzzzzzz/hot/crawler/wsj"
	"github.com/chamzzzzzz/hot/crawler/xueqiu"
	"github.com/chamzzzzzz/hot/crawler/yiche"
	"github.com/chamzzzzzz/hot/crawler/yystv"
	"github.com/chamzzzzzz/hot/crawler/zhiguf"
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
	hc.crawlers = append(hc.crawlers, &ithome.Crawler{})
	hc.crawlers = append(hc.crawlers, &ithome.Crawler{ithome.BoardGame})
	hc.crawlers = append(hc.crawlers, &thepaper.Crawler{})
	hc.crawlers = append(hc.crawlers, &kr36.Crawler{})
	hc.crawlers = append(hc.crawlers, &bilibili.Crawler{})
	hc.crawlers = append(hc.crawlers, &netease.Crawler{})
	hc.crawlers = append(hc.crawlers, &thecover.Crawler{})
	hc.crawlers = append(hc.crawlers, &wsj.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &techweb.Crawler{})
	hc.crawlers = append(hc.crawlers, &infoq.Crawler{})
	hc.crawlers = append(hc.crawlers, &pojie52.Crawler{})
	hc.crawlers = append(hc.crawlers, &daniu.Crawler{})
	hc.crawlers = append(hc.crawlers, &sogou.Crawler{})
	hc.crawlers = append(hc.crawlers, &sogou.Crawler{sogou.Weixin})
	hc.crawlers = append(hc.crawlers, &sogou.Crawler{sogou.Baike})
	hc.crawlers = append(hc.crawlers, &douban.Crawler{})
	hc.crawlers = append(hc.crawlers, &douban.Crawler{douban.Movie})
	hc.crawlers = append(hc.crawlers, &hupu.Crawler{})
	hc.crawlers = append(hc.crawlers, &hupu.Crawler{hupu.Basketball})
	hc.crawlers = append(hc.crawlers, &hupu.Crawler{hupu.Football})
	hc.crawlers = append(hc.crawlers, &hupu.Crawler{hupu.Gambia})
	hc.crawlers = append(hc.crawlers, &chinaz.Crawler{})
	hc.crawlers = append(hc.crawlers, &kanxue.Crawler{})
	hc.crawlers = append(hc.crawlers, &kanxue.Crawler{kanxue.BBS})
	hc.crawlers = append(hc.crawlers, &cnbeta.Crawler{})
	hc.crawlers = append(hc.crawlers, &so360.Crawler{})
	hc.crawlers = append(hc.crawlers, &haokan.Crawler{})
	hc.crawlers = append(hc.crawlers, &nowcoder.Crawler{})
	hc.crawlers = append(hc.crawlers, &toutiaoio.Crawler{})
	hc.crawlers = append(hc.crawlers, &oschina.Crawler{})
	hc.crawlers = append(hc.crawlers, &cto51.Crawler{})
	hc.crawlers = append(hc.crawlers, &gameres.Crawler{})
	hc.crawlers = append(hc.crawlers, &investing.Crawler{})
	hc.crawlers = append(hc.crawlers, &gitchat.Crawler{})
	hc.crawlers = append(hc.crawlers, &cls.Crawler{})
	hc.crawlers = append(hc.crawlers, &jin10.Crawler{})
	hc.crawlers = append(hc.crawlers, &jqka10.Crawler{})
	hc.crawlers = append(hc.crawlers, &csdn.Crawler{})
	hc.crawlers = append(hc.crawlers, &xueqiu.Crawler{})
	hc.crawlers = append(hc.crawlers, &eastmoney.Crawler{})
	hc.crawlers = append(hc.crawlers, &sohu.Crawler{})
	hc.crawlers = append(hc.crawlers, &donews.Crawler{})
	hc.crawlers = append(hc.crawlers, &kuaishou.Crawler{})
	hc.crawlers = append(hc.crawlers, &ifeng.Crawler{})
	hc.crawlers = append(hc.crawlers, &yystv.Crawler{})
	hc.crawlers = append(hc.crawlers, &globaltimes.Crawler{})
	hc.crawlers = append(hc.crawlers, &bjnews.Crawler{})
	hc.crawlers = append(hc.crawlers, &rfa.Crawler{rfa.Mandarin, proxy})
	hc.crawlers = append(hc.crawlers, &rfa.Crawler{rfa.Cantonese, proxy})
	hc.crawlers = append(hc.crawlers, &rfa.Crawler{rfa.English, proxy})
	hc.crawlers = append(hc.crawlers, &tianya.Crawler{})
	hc.crawlers = append(hc.crawlers, &cnblogs.Crawler{})
	hc.crawlers = append(hc.crawlers, &jisilu.Crawler{})
	hc.crawlers = append(hc.crawlers, &mydrivers.Crawler{})
	hc.crawlers = append(hc.crawlers, &odaily.Crawler{})
	hc.crawlers = append(hc.crawlers, &readhub.Crawler{})
	hc.crawlers = append(hc.crawlers, &chinanews.Crawler{})
	hc.crawlers = append(hc.crawlers, &panews.Crawler{})
	hc.crawlers = append(hc.crawlers, &zhiguf.Crawler{})
	hc.crawlers = append(hc.crawlers, &btc8.Crawler{})
	hc.crawlers = append(hc.crawlers, &jinse.Crawler{})
	hc.crawlers = append(hc.crawlers, &dongchedi.Crawler{})
	hc.crawlers = append(hc.crawlers, &leikeji.Crawler{})
	hc.crawlers = append(hc.crawlers, &chiphell.Crawler{})
	hc.crawlers = append(hc.crawlers, &cninfo.Crawler{})
	hc.crawlers = append(hc.crawlers, &hibor.Crawler{})
	hc.crawlers = append(hc.crawlers, &gelonghui.Crawler{})
	hc.crawlers = append(hc.crawlers, &fxbaogao.Crawler{})
	hc.crawlers = append(hc.crawlers, &nytimes.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &bbc.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &yiche.Crawler{})
	hc.crawlers = append(hc.crawlers, &gamersky.Crawler{})
	hc.crawlers = append(hc.crawlers, &solidot.Crawler{})
	hc.crawlers = append(hc.crawlers, &pearvideo.Crawler{})
	hc.crawlers = append(hc.crawlers, &sspai.Crawler{})
	hc.crawlers = append(hc.crawlers, &timecom.Crawler{proxy})

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
