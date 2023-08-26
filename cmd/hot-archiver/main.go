package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"mime"
	"net"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/archiver/database"
	"github.com/chamzzzzzz/hot/archiver/file"
	"github.com/chamzzzzzz/hot/crawler"
	_ "github.com/go-sql-driver/mysql"
)

var (
	proxy     = os.Getenv("HOT_ARCHIVER_PROXY")
	board     = os.Getenv("HOT_ARCHIVER_BOARD")
	mode      = os.Getenv("HOT_ARCHIVER_MODE")
	once      = os.Getenv("HOT_ARCHIVER_ONCE")
	dn        = os.Getenv("HOT_ARCHIVER_DATABASE_DRIVER_NAME")
	dsn       = os.Getenv("HOT_ARCHIVER_DATABASE_DATA_SOURCE_NAME")
	archivers []hot.Archiver
	crawlers  []*crawler.Crawler
	boards    = map[string][]string{
		"china-popular": {"baidu", "weibo", "toutiao", "douyin", "kuaishou", "bilibili"},
		"global": {
			"cnbeta", "economist", "ft", "ftchinese", "github", "hket", "kyodonews", "nytimes",
			"rfa", "thehill", "timecom", "v2ex", "voacantonese", "voachinese", "wikipedia", "wsj",
		},
	}
)

var (
	addr   = os.Getenv("HOT_ARCHIVER_SMTP_ADDR")
	user   = os.Getenv("HOT_ARCHIVER_SMTP_USER")
	pass   = os.Getenv("HOT_ARCHIVER_SMTP_PASS")
	source = "From: {{.From}}\r\nTo: {{.To}}\r\nSubject: {{.Subject}}\r\n\r\n{{.Body}}"
	tpl    *template.Template
	stats  = make(map[string]*stat)
)

type stat struct {
	Count             int
	CrawlFailCount    int
	CrawlNothingCount int
	ArchiveFailCount  int
}

func main() {
	flag.StringVar(&proxy, "proxy", proxy, "proxy url")
	flag.StringVar(&board, "board", board, "china-popular(default), china, global, all, or custom comma separated driver names")
	flag.StringVar(&mode, "mode", mode, "file(default), database, all")
	flag.StringVar(&once, "once", once, "archive one time")
	flag.StringVar(&dn, "dn", dn, "database driver name")
	flag.StringVar(&dsn, "dsn", dsn, "database data source name")
	flag.StringVar(&addr, "addr", addr, "notification smtp addr")
	flag.StringVar(&user, "user", user, "notification smtp user")
	flag.StringVar(&pass, "pass", pass, "notification smtp pass")
	flag.Parse()

	funcs := template.FuncMap{
		"bencoding": mime.BEncoding.Encode,
	}
	tpl = template.Must(template.New("mail").Funcs(funcs).Parse(source))

	if mode == "" {
		mode = "file"
	}
	if mode == "all" {
		mode = "file,database"
	}
	for _, m := range strings.Split(mode, ",") {
		switch m {
		case "database":
			if dn == "" {
				log.Println("database driver name is empty")
				return
			}
			if dsn == "" {
				log.Println("database data source name is empty")
				return
			}
			a := &database.Archiver{DriverName: dn, DataSourceName: dsn}
			archivers = append(archivers, a)
		case "file":
			a := &file.Archiver{}
			archivers = append(archivers, a)
		default:
			log.Printf("unknown mode %s\n", mode)
			return
		}
	}

	board, drivers := parse(board)
	log.Printf("proxy=%s\n", proxy)
	log.Printf("board=%s\n", board)
	for _, driver := range drivers {
		log.Printf("driver=%s\n", driver)
	}
	for _, archiver := range archivers {
		log.Printf("archiver=%s\n", archiver.Name())
	}
	for _, driverName := range drivers {
		c, err := crawler.Open(crawler.Option{DriverName: driverName, Proxy: proxy})
		if err != nil {
			log.Printf("[%s] open crawler failed, err=%s\n", driverName, err)
			return
		}
		crawlers = append(crawlers, c)
		stats[driverName] = &stat{}
	}
	for {
		archive()
		if once == "1" {
			break
		}
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		log.Printf("next archive at %s\n", next.Format("2006-01-02 15:04:05"))
		time.Sleep(next.Sub(now))
	}
}

func parse(board string) (string, []string) {
	var drivers []string
	if board == "" {
		board = "china-popular"
	}
	if board == "all" {
		drivers = crawler.Drivers()
	} else if board == "china" {
		m := make(map[string]bool)
		for _, v := range boards["global"] {
			m[v] = true
		}
		for _, name := range crawler.Drivers() {
			if !m[name] {
				drivers = append(drivers, name)
			}
		}
	} else {
		_driver, ok := boards[board]
		if !ok {
			drivers = strings.Split(board, ",")
			board = "custom"
		} else {
			drivers = _driver
		}
	}
	return board, drivers
}

func archive() {
	log.Printf("start archive at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	t := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, c := range crawlers {
		wg.Add(1)
		go func(c *crawler.Crawler) {
			defer wg.Done()
			mu.Lock()
			stats[c.Name()].Count++
			mu.Unlock()
			board, err := c.Crawl()
			if err != nil {
				for i := 0; i < 4; i++ {
					log.Printf("[%s] crawl failed, will retry(%d) after 5 seconds. err=%s\n", c.Name(), i+1, err)
					time.Sleep(5 * time.Second)
					board, err = c.Crawl()
					if err == nil {
						break
					}
				}
			}
			if err != nil {
				log.Printf("[%s] crawl failed, err=%s\n", c.Name(), err)
				mu.Lock()
				stats[c.Name()].CrawlFailCount++
				mu.Unlock()
				return
			}
			if len(board.Hots) == 0 {
				log.Printf("[%s] crawl nothing\n", c.Name())
				mu.Lock()
				stats[c.Name()].CrawlNothingCount++
				mu.Unlock()
				return
			}
			for _, archiver := range archivers {
				_, err = archiver.Archive(board)
				if err != nil {
					log.Printf("[%s] archive failed, err=%s\n", c.Name(), err)
					mu.Lock()
					stats[c.Name()].ArchiveFailCount++
					mu.Unlock()
					return
				}
			}
		}(c)
	}
	wg.Wait()
	log.Printf("archive used %v\n", time.Since(t))
	log.Printf("finish archive at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	notification()
}

func notification() {
	type Data struct {
		From    string
		To      string
		Subject string
		Body    string
	}

	if addr == "" {
		log.Printf("send notification skip. addr is empty\n")
		return
	}

	if once == "" {
		now := time.Now()
		if now.Hour() != 19 {
			log.Printf("send notification skip. time is not 19:00\n")
			return
		}
	}

	var failed = make(map[int]string)
	for _, c := range crawlers {
		s := stats[c.Name()]
		if s.CrawlFailCount > 0 {
			failed[1] += fmt.Sprintf("%s (%d/%d)\n", c.Name(), s.CrawlFailCount, s.Count)
		}
		if s.CrawlNothingCount > 0 {
			failed[2] += fmt.Sprintf("%s (%d/%d)\n", c.Name(), s.CrawlNothingCount, s.Count)
		}
		if s.ArchiveFailCount > 0 {
			failed[3] += fmt.Sprintf("%s (%d/%d)\n", c.Name(), s.ArchiveFailCount, s.Count)
		}
		s.Count = 0
		s.CrawlFailCount = 0
		s.CrawlNothingCount = 0
		s.ArchiveFailCount = 0
	}
	if len(failed) == 0 {
		log.Printf("send notification skip. no failed\n")
		return
	}

	subject := "「HA」异常发生"
	body := ""
	for i := 1; i <= 3; i++ {
		if failed[i] != "" {
			desc := ""
			switch i {
			case 1:
				desc = "获取失败的榜单："
			case 2:
				desc = "获取空白的榜单："
			case 3:
				desc = "归档失败的榜单："
			}
			body += fmt.Sprintf("%s\n%s\n", desc, failed[i])
		}
	}

	log.Printf("sending notification...")
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Printf("send notification fail. err='%s'\n", err)
		return
	}

	data := Data{
		From:    fmt.Sprintf("%s <%s>", mime.BEncoding.Encode("UTF-8", "Monitor"), user),
		To:      user,
		Subject: mime.BEncoding.Encode("UTF-8", subject),
		Body:    body,
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		log.Printf("send notification fail. err='%s'\n", err)
		return
	}

	auth := smtp.PlainAuth("", user, pass, host)
	if err := smtp.SendMail(addr, auth, user, []string{user}, buf.Bytes()); err != nil {
		log.Printf("send notification fail. err='%s'\n", err)
	}
	log.Printf("send notification success.\n")
}
