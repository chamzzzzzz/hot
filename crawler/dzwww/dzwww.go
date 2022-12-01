package dzwww

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
)

const (
	HotSearch = "hotsearch"
	HotNews   = "hotnews"
)

var act = map[string]string{
	HotSearch: "getPCHotSearch",
	HotNews:   "getPCHotNews",
}

type Crawler struct {
	Catalog string
}

func (c *Crawler) Name() string {
	return "dzwww"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.Catalog {
	case HotSearch, HotNews:
		return c.withCatalog(c.Catalog, nil)
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	for _, catalog := range []string{HotSearch, HotNews} {
		if _, err := c.withCatalog(catalog, board); err != nil {
			return nil, err
		}
	}
	return board, nil
}

func (c *Crawler) withCatalog(catalog string, board *hot.Board) (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://w.dzwww.com/?act=%s", act[catalog]), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.Set("Referer", "https://w.dzwww.com/")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var body []body
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	for _, data := range body {
		board.AppendTitleURLCatalog(data.Title, "https:"+data.RealURL, catalog)
	}
	return board, nil
}

type body struct {
	Title   string `json:"title"`
	RealURL string `json:"real_url"`
}
