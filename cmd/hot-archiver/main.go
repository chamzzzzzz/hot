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

	"github.com/chamzzzzzz/hot/archiver/file"
	"github.com/chamzzzzzz/hot/crawler"
)

var (
	proxy    = os.Getenv("HOT_ARCHIVER_PROXY")
	board    = os.Getenv("HOT_ARCHIVER_BOARD")
	archiver = &file.Archiver{}
	crawlers []*crawler.Crawler
	boards   = map[string][]string{
		"china":  {"baidu", "weibo", "toutiao", "douyin", "kuaishou", "bilibili"},
		"global": {"wsj", "bbc"},
	}
)

var (
	addr   = os.Getenv("HOT_ARCHIVER_SMTP_ADDR")
	user   = os.Getenv("HOT_ARCHIVER_SMTP_USER")
	pass   = os.Getenv("HOT_ARCHIVER_SMTP_PASS")
	source = "From: {{.From}}\r\nTo: {{.To}}\r\nSubject: {{.Subject}}\r\n\r\n{{.Body}}"
	tpl    *template.Template
)

func main() {
	flag.StringVar(&proxy, "proxy", proxy, "proxy url")
	flag.StringVar(&board, "board", board, "china, global, all, or custom comma separated driver names")
	flag.StringVar(&addr, "addr", addr, "notification smtp addr")
	flag.StringVar(&user, "user", user, "notification smtp user")
	flag.StringVar(&pass, "pass", pass, "notification smtp pass")
	flag.Parse()

	funcs := template.FuncMap{
		"bencoding": mime.BEncoding.Encode,
	}
	tpl = template.Must(template.New("mail").Funcs(funcs).Parse(source))

	board, drivers := parse(board)
	log.Printf("proxy=%s\n", proxy)
	log.Printf("board=%s\n", board)
	log.Printf("drivers=%v\n", drivers)
	log.Printf("archiver=%s\n", archiver.Name())
	for _, driverName := range drivers {
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

func parse(board string) (string, []string) {
	var drivers []string
	if board == "" {
		board = "china"
	}
	if board == "all" {
		drivers = crawler.Drivers()
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
	var errmsgs []string
	for _, c := range crawlers {
		wg.Add(1)
		go func(c *crawler.Crawler) {
			defer wg.Done()
			board, err := c.Crawl()
			if err != nil {
				errmsg := fmt.Sprintf("[%s] crawl failed, err=%s\n", c.Name(), err)
				log.Print(errmsg)
				mu.Lock()
				errmsgs = append(errmsgs, errmsg)
				mu.Unlock()
				return
			}
			archiver.Archive(board)
		}(c)
	}
	wg.Wait()
	log.Printf("archive used %v\n", time.Since(t))
	log.Printf("finish archive at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	if len(errmsgs) > 0 {
		notification("「HA」异常发生", strings.Join(errmsgs, ""))
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
		From:    fmt.Sprintf("%s <%s>", mime.BEncoding.Encode("UTF-8", "HA Monitor"), user),
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
