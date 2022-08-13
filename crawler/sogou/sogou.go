package sogou

import (
	"encoding/json"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	Weixin = "weixin"
	Baike  = "baike"
)

type Crawler struct {
	BoardName string
}

func (c *Crawler) Name() string {
	switch c.BoardName {
	case Weixin:
		return "sogou_x_weixin"
	case Baike:
		return "sogou_x_baike"
	default:
		return "sogou"
	}
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.BoardName {
	case Weixin:
		return c.sogou_x_weixin()
	case Baike:
		return c.sogou_x_baike()
	default:
		return c.sogou()
	}
}

func (c *Crawler) sogou_x_weixin() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://weixin.sogou.com", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(html))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	ol := dom.Find("ol", "class", "hot-news")
	if ol.Error != nil {
		return nil, ol.Error
	}
	for _, a := range ol.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}

func (c *Crawler) sogou_x_baike() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://baike.sogou.com", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(html))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	ol := dom.Find("ol", "class", "hot_rank")
	if ol.Error != nil {
		return nil, ol.Error
	}
	for _, a := range ol.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}

func (c *Crawler) sogou() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.sogou.com/suggnew/hotwords", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyJson := []string{}
	if err := json.Unmarshal(body[20:], &bodyJson); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, word := range bodyJson {
		board.Append(word, "", date)
	}
	return board, nil
}
