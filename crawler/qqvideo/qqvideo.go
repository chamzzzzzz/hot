package qqvideo

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	General = "rank"
	TV      = "tv"
	Variety = "variety"
	Cartoon = "cartoon"
	Child   = "child"
	Movie   = "movie"
	Doco    = "doco"
	Games   = "games"
	Music   = "music"
)

type Crawler struct {
	BoardName string
}

func (c *Crawler) Name() string {
	switch c.BoardName {
	case General:
		return "qqvideo_x_general"
	case TV, Variety, Cartoon, Child, Movie, Doco, Games, Music:
		return "qqvideo_x_" + c.BoardName
	default:
		return "qqvideo"
	}
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.BoardName {
	case General, TV, Variety, Cartoon, Child, Movie, Doco, Games, Music:
		return c.qqvideoWithChannel(c.BoardName)
	default:
		return c.qqvideo()
	}
}

func (c *Crawler) qqvideo() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://v.qq.com/biu/ranks/?t=hotsearch", nil)
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
	for _, li := range dom.FindAllStrict("li", "class", "item item_odd item_1") {
		a := li.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		title := strings.TrimSpace(a.Attrs()["title"])
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}

func (c *Crawler) qqvideoWithChannel(channel string) (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://v.qq.com/biu/ranks/?t=hotsearch&channel=%s", channel), nil)
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
	for _, div := range dom.FindAllStrict("div", "class", "item item_a") {
		a := div.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		title := strings.TrimSpace(a.Text())
		summary := a.Attrs()["href"]
		board.Append(title, summary, date)
	}
	return board, nil
}
