package jrj

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	Tech  = "tech"
	House = "house"
)

type Crawler struct {
	BoardName string
}

func (c *Crawler) Name() string {
	switch c.BoardName {
	case Tech:
		return "jrj_x_tech"
	case House:
		return "jrj_x_house"
	default:
		return "jrj"
	}
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.BoardName {
	case Tech:
		return c.jrj_x_tech()
	case House:
		return c.jrj_x_house()
	default:
		return c.jrj()
	}
}

func (c *Crawler) jrj_x_tech() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://tech.jrj.com.cn", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(html))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	div := dom.FindStrict("div", "class", "hotart hotnews")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}

func (c *Crawler) jrj_x_house() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://house.jrj.com.cn", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(html))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	div := dom.FindStrict("div", "class", "hotart")
	if div.Error != nil {
		return nil, div.Error
	}
	for _, a := range div.FindAllStrict("a") {
		title := strings.TrimSpace(a.Text())
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}

func (c *Crawler) jrj() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://finance.jrj.com.cn/list/industrynews.shtml", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(html))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, id := range []string{"con_1", "con_2"} {
		ul := dom.Find("ul", "id", id)
		if ul.Error != nil {
			return nil, ul.Error
		}
		for _, a := range ul.FindAll("a") {
			title := strings.TrimSpace(a.Text())
			summary := a.Attrs()["href"]
			board.Append(title, summary, date)
		}
	}
	return board, nil
}
