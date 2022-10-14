package kanxue

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	News = "news"
	BBS  = "bbs"
)

type Crawler struct {
	Catalog string
}

func (c *Crawler) Name() string {
	return "kanxue"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.Catalog {
	case News:
		return c.news()
	case BBS:
		return c.bbs()
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	b1, err := c.news()
	if err != nil {
		return nil, err
	}
	b2, err := c.bbs()
	if err != nil {
		return nil, err
	}
	for _, hot := range b1.Hots {
		board.Append0(hot)
	}
	for _, hot := range b2.Hots {
		board.Append0(hot)
	}
	return board, nil
}

func (c *Crawler) bbs() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://bbs.pediy.com/thread-hotlist-all-0.htm", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	tbody := dom.Find("tbody", "id", "arctilelist")
	if tbody.Error != nil {
		return nil, tbody.Error
	}
	for _, tr := range tbody.FindAll("tr") {
		a := tr.Find("a", "class", "text-white")
		if a.Error != nil {
			return nil, a.Error
		}
		title := strings.TrimSpace(a.Text())
		url := "https://bbs.pediy.com/" + strings.TrimSpace(a.Attrs()["href"])
		board.Append4(title, "", url, BBS)
	}
	return board, nil
}

func (c *Crawler) news() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.kanxue.com", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	for _, div := range dom.FindAllStrict("div", "class", "pr-3 pb-2 mb-2 position-relative") {
		a := div.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Attrs()["href"])
		board.Append4(title, "", url, News)
	}
	return board, nil
}
