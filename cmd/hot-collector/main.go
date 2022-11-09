package main

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/archiver/database"
	"github.com/chamzzzzzz/hot/crawler/acfun"
	"github.com/chamzzzzzz/hot/crawler/baidu"
	"github.com/chamzzzzzz/hot/crawler/baijing"
	"github.com/chamzzzzzz/hot/crawler/bbc"
	"github.com/chamzzzzzz/hot/crawler/bilibili"
	"github.com/chamzzzzzz/hot/crawler/bjnews"
	"github.com/chamzzzzzz/hot/crawler/btc8"
	"github.com/chamzzzzzz/hot/crawler/chinanews"
	"github.com/chamzzzzzz/hot/crawler/chinaz"
	"github.com/chamzzzzzz/hot/crawler/chiphell"
	"github.com/chamzzzzzz/hot/crawler/chunfengxing"
	"github.com/chamzzzzzz/hot/crawler/cls"
	"github.com/chamzzzzzz/hot/crawler/cnbeta"
	"github.com/chamzzzzzz/hot/crawler/cnblogs"
	"github.com/chamzzzzzz/hot/crawler/cninfo"
	"github.com/chamzzzzzz/hot/crawler/credit51"
	"github.com/chamzzzzzz/hot/crawler/csdn"
	"github.com/chamzzzzzz/hot/crawler/cto51"
	"github.com/chamzzzzzz/hot/crawler/ctoutiao"
	"github.com/chamzzzzzz/hot/crawler/cyzone"
	"github.com/chamzzzzzz/hot/crawler/daniu"
	"github.com/chamzzzzzz/hot/crawler/donews"
	"github.com/chamzzzzzz/hot/crawler/dongchedi"
	"github.com/chamzzzzzz/hot/crawler/dongqiudi"
	"github.com/chamzzzzzz/hot/crawler/douban"
	"github.com/chamzzzzzz/hot/crawler/douyin"
	"github.com/chamzzzzzz/hot/crawler/dxy"
	"github.com/chamzzzzzz/hot/crawler/eastmoney"
	"github.com/chamzzzzzz/hot/crawler/economist"
	"github.com/chamzzzzzz/hot/crawler/eeo"
	"github.com/chamzzzzzz/hot/crawler/fortunechina"
	"github.com/chamzzzzzz/hot/crawler/ft"
	"github.com/chamzzzzzz/hot/crawler/ftchinese"
	"github.com/chamzzzzzz/hot/crawler/futu"
	"github.com/chamzzzzzz/hot/crawler/fxbaogao"
	"github.com/chamzzzzzz/hot/crawler/gameres"
	"github.com/chamzzzzzz/hot/crawler/gamersky"
	"github.com/chamzzzzzz/hot/crawler/gelonghui"
	"github.com/chamzzzzzz/hot/crawler/github"
	"github.com/chamzzzzzz/hot/crawler/globaltimes"
	"github.com/chamzzzzzz/hot/crawler/guancha"
	"github.com/chamzzzzzz/hot/crawler/haokan"
	"github.com/chamzzzzzz/hot/crawler/hibor"
	"github.com/chamzzzzzz/hot/crawler/hupu"
	"github.com/chamzzzzzz/hot/crawler/huxiu"
	"github.com/chamzzzzzz/hot/crawler/ifeng"
	"github.com/chamzzzzzz/hot/crawler/igao7"
	"github.com/chamzzzzzz/hot/crawler/im2maker"
	"github.com/chamzzzzzz/hot/crawler/infoq"
	"github.com/chamzzzzzz/hot/crawler/infoqcom"
	"github.com/chamzzzzzz/hot/crawler/investing"
	"github.com/chamzzzzzz/hot/crawler/iqiyi"
	"github.com/chamzzzzzz/hot/crawler/iresearch"
	"github.com/chamzzzzzz/hot/crawler/ithome"
	"github.com/chamzzzzzz/hot/crawler/jiemian"
	"github.com/chamzzzzzz/hot/crawler/jin10"
	"github.com/chamzzzzzz/hot/crawler/jinse"
	"github.com/chamzzzzzz/hot/crawler/jisilu"
	"github.com/chamzzzzzz/hot/crawler/jqka10"
	"github.com/chamzzzzzz/hot/crawler/jrj"
	"github.com/chamzzzzzz/hot/crawler/kanxue"
	"github.com/chamzzzzzz/hot/crawler/kr36"
	"github.com/chamzzzzzz/hot/crawler/kuaishou"
	"github.com/chamzzzzzz/hot/crawler/kugou"
	"github.com/chamzzzzzz/hot/crawler/kyodonews"
	"github.com/chamzzzzzz/hot/crawler/lanjinger"
	"github.com/chamzzzzzz/hot/crawler/leikeji"
	"github.com/chamzzzzzz/hot/crawler/maoyan"
	"github.com/chamzzzzzz/hot/crawler/mydrivers"
	"github.com/chamzzzzzz/hot/crawler/netease"
	"github.com/chamzzzzzz/hot/crawler/nowcoder"
	"github.com/chamzzzzzz/hot/crawler/nytimes"
	"github.com/chamzzzzzz/hot/crawler/odaily"
	"github.com/chamzzzzzz/hot/crawler/oschina"
	"github.com/chamzzzzzz/hot/crawler/panews"
	"github.com/chamzzzzzz/hot/crawler/pearvideo"
	"github.com/chamzzzzzz/hot/crawler/pojie52"
	"github.com/chamzzzzzz/hot/crawler/pudn"
	"github.com/chamzzzzzz/hot/crawler/qqnews"
	"github.com/chamzzzzzz/hot/crawler/qqvideo"
	"github.com/chamzzzzzz/hot/crawler/readhub"
	"github.com/chamzzzzzz/hot/crawler/rfa"
	"github.com/chamzzzzzz/hot/crawler/rt"
	"github.com/chamzzzzzz/hot/crawler/semiunion"
	"github.com/chamzzzzzz/hot/crawler/so360"
	"github.com/chamzzzzzz/hot/crawler/sogou"
	"github.com/chamzzzzzz/hot/crawler/sohu"
	"github.com/chamzzzzzz/hot/crawler/solidot"
	"github.com/chamzzzzzz/hot/crawler/sputniknews"
	"github.com/chamzzzzzz/hot/crawler/sspai"
	"github.com/chamzzzzzz/hot/crawler/takungpao"
	"github.com/chamzzzzzz/hot/crawler/taptap"
	"github.com/chamzzzzzz/hot/crawler/techweb"
	"github.com/chamzzzzzz/hot/crawler/thecover"
	"github.com/chamzzzzzz/hot/crawler/theguardian"
	"github.com/chamzzzzzz/hot/crawler/thehill"
	"github.com/chamzzzzzz/hot/crawler/thepaper"
	"github.com/chamzzzzzz/hot/crawler/tianya"
	"github.com/chamzzzzzz/hot/crawler/tieba"
	"github.com/chamzzzzzz/hot/crawler/timecom"
	"github.com/chamzzzzzz/hot/crawler/toutiao"
	"github.com/chamzzzzzz/hot/crawler/toutiaoio"
	"github.com/chamzzzzzz/hot/crawler/v2ex"
	"github.com/chamzzzzzz/hot/crawler/vrtuoluo"
	"github.com/chamzzzzzz/hot/crawler/wallstreetcn"
	"github.com/chamzzzzzz/hot/crawler/weibo"
	"github.com/chamzzzzzz/hot/crawler/wikipedia"
	"github.com/chamzzzzzz/hot/crawler/wsj"
	"github.com/chamzzzzzz/hot/crawler/xueqiu"
	"github.com/chamzzzzzz/hot/crawler/yfchuhai"
	"github.com/chamzzzzzz/hot/crawler/yicai"
	"github.com/chamzzzzzz/hot/crawler/yiche"
	"github.com/chamzzzzzz/hot/crawler/youxituoluo"
	"github.com/chamzzzzzz/hot/crawler/yystv"
	"github.com/chamzzzzzz/hot/crawler/zaker"
	"github.com/chamzzzzzz/hot/crawler/zaobao"
	"github.com/chamzzzzzz/hot/crawler/zhiguf"
	"github.com/chamzzzzzz/hot/crawler/zhihu"
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
	hc.crawlers = append(hc.crawlers, &douban.Crawler{})
	hc.crawlers = append(hc.crawlers, &hupu.Crawler{})
	hc.crawlers = append(hc.crawlers, &chinaz.Crawler{})
	hc.crawlers = append(hc.crawlers, &kanxue.Crawler{})
	hc.crawlers = append(hc.crawlers, &cnbeta.Crawler{})
	hc.crawlers = append(hc.crawlers, &so360.Crawler{})
	hc.crawlers = append(hc.crawlers, &haokan.Crawler{})
	hc.crawlers = append(hc.crawlers, &nowcoder.Crawler{})
	hc.crawlers = append(hc.crawlers, &toutiaoio.Crawler{})
	hc.crawlers = append(hc.crawlers, &oschina.Crawler{})
	hc.crawlers = append(hc.crawlers, &cto51.Crawler{})
	hc.crawlers = append(hc.crawlers, &gameres.Crawler{})
	hc.crawlers = append(hc.crawlers, &investing.Crawler{})
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
	hc.crawlers = append(hc.crawlers, &rfa.Crawler{Proxy: proxy})
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
	hc.crawlers = append(hc.crawlers, &kyodonews.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &huxiu.Crawler{})
	hc.crawlers = append(hc.crawlers, &takungpao.Crawler{})
	hc.crawlers = append(hc.crawlers, &dongqiudi.Crawler{})
	hc.crawlers = append(hc.crawlers, &zaker.Crawler{})
	hc.crawlers = append(hc.crawlers, &qqnews.Crawler{})
	hc.crawlers = append(hc.crawlers, &ft.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &ftchinese.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &theguardian.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &futu.Crawler{})
	hc.crawlers = append(hc.crawlers, &wallstreetcn.Crawler{})
	hc.crawlers = append(hc.crawlers, &iresearch.Crawler{})
	hc.crawlers = append(hc.crawlers, &cyzone.Crawler{})
	hc.crawlers = append(hc.crawlers, &zaobao.Crawler{})
	hc.crawlers = append(hc.crawlers, &eeo.Crawler{})
	hc.crawlers = append(hc.crawlers, &semiunion.Crawler{})
	hc.crawlers = append(hc.crawlers, &igao7.Crawler{})
	hc.crawlers = append(hc.crawlers, &chunfengxing.Crawler{})
	hc.crawlers = append(hc.crawlers, &credit51.Crawler{})
	hc.crawlers = append(hc.crawlers, &qqvideo.Crawler{})
	hc.crawlers = append(hc.crawlers, &jrj.Crawler{})
	hc.crawlers = append(hc.crawlers, &jiemian.Crawler{})
	hc.crawlers = append(hc.crawlers, &lanjinger.Crawler{})
	hc.crawlers = append(hc.crawlers, &fortunechina.Crawler{})
	hc.crawlers = append(hc.crawlers, &guancha.Crawler{})
	hc.crawlers = append(hc.crawlers, &economist.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &wikipedia.Crawler{proxy})
	hc.crawlers = append(hc.crawlers, &sputniknews.Crawler{})
	hc.crawlers = append(hc.crawlers, &rt.Crawler{})
	hc.crawlers = append(hc.crawlers, &thehill.Crawler{})
	hc.crawlers = append(hc.crawlers, &yicai.Crawler{})
	hc.crawlers = append(hc.crawlers, &iqiyi.Crawler{})
	hc.crawlers = append(hc.crawlers, &acfun.Crawler{})
	hc.crawlers = append(hc.crawlers, &infoqcom.Crawler{})
	hc.crawlers = append(hc.crawlers, &youxituoluo.Crawler{})
	hc.crawlers = append(hc.crawlers, &vrtuoluo.Crawler{})
	hc.crawlers = append(hc.crawlers, &im2maker.Crawler{})
	hc.crawlers = append(hc.crawlers, &ctoutiao.Crawler{})
	hc.crawlers = append(hc.crawlers, &baijing.Crawler{})
	hc.crawlers = append(hc.crawlers, &yfchuhai.Crawler{})
	hc.crawlers = append(hc.crawlers, &dxy.Crawler{})
	hc.crawlers = append(hc.crawlers, &kugou.Crawler{})
	hc.crawlers = append(hc.crawlers, &taptap.Crawler{})
	hc.crawlers = append(hc.crawlers, &maoyan.Crawler{})
	hc.crawlers = append(hc.crawlers, &pudn.Crawler{})

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
