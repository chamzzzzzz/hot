package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net"
	"net/smtp"
	"os"
	"sync"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/archiver/database"
	"github.com/chamzzzzzz/hot/crawler"
	"github.com/robfig/cron/v3"
)

var (
	logger = log.New(os.Stdout, "collector: ", log.Ldate|log.Lmicroseconds)
)

var (
	addr   = os.Getenv("HOT_COLLECTOR_SMTP_ADDR")
	user   = os.Getenv("HOT_COLLECTOR_SMTP_USER")
	pass   = os.Getenv("HOT_COLLECTOR_SMTP_PASS")
	source = "From: {{.From}}\r\nTo: {{.To}}\r\nSubject: {{.Subject}}\r\n\r\n{{.Body}}"
	tpl    *template.Template
)

type HotCollector struct {
	option   crawler.Option
	crawlers []*crawler.Crawler
	archiver hot.Archiver
}

func (hc *HotCollector) Start() error {
	flag.StringVar(&addr, "addr", addr, "notification smtp addr")
	flag.StringVar(&user, "user", user, "notification smtp user")
	flag.StringVar(&pass, "pass", pass, "notification smtp pass")
	flag.Parse()
	funcs := template.FuncMap{
		"bencoding": mime.BEncoding.Encode,
	}
	tpl = template.Must(template.New("mail").Funcs(funcs).Parse(source))

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
	var mu sync.Mutex
	failed := make(map[int][]string)
	for _, crawler := range hc.crawlers {
		wg.Add(1)
		go func(crawler hot.Crawler) {
			defer wg.Done()
			board, err := crawler.Crawl()
			if err != nil {
				logger.Printf("crawl, error='%s', crawler=%s\n", err, crawler.Name())
				mu.Lock()
				failed[1] = append(failed[1], crawler.Name())
				mu.Unlock()
				return
			}
			logger.Printf("crawl, crawler=%s, board=%s, count=%d\n", crawler.Name(), board.Name, len(board.Hots))

			var archived = 0
			if archived, err = hc.archiver.Archive(board); err != nil {
				logger.Printf("archive, error='%s', archiver=%s, board=%s, count=%d, archived=%d\n", err, hc.archiver.Name(), board.Name, len(board.Hots), archived)
				mu.Lock()
				failed[2] = append(failed[2], crawler.Name())
				mu.Unlock()
				return
			}
			logger.Printf("archive, archiver=%s, board=%s, count=%d, arvhiced=%d\n", hc.archiver.Name(), board.Name, len(board.Hots), archived)
		}(crawler)
	}
	wg.Wait()
	logger.Printf("run, used=%v", time.Since(t1))
	if len(failed) > 0 {
		body := ""
		if len(failed[1]) > 0 {
			body += "获取失败的榜单：\n"
			for _, name := range failed[1] {
				body += fmt.Sprintf("%s\n", name)
			}
			body += "\n"
		}
		if len(failed[2]) > 0 {
			body += "归档失败的榜单：\n"
			for _, name := range failed[2] {
				body += fmt.Sprintf("%s\n", name)
			}
			body += "\n"
		}
		notification("「HC」异常发生", body)
	}
}

func notification(subject, body string) {
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

func main() {
	hc := &HotCollector{}
	err := hc.Start()
	if err != nil {
		logger.Printf("start, error='%s'\n", err)
		os.Exit(1)
	}
}
