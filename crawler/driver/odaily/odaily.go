package odaily

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	DriverName = "odaily"
)

type Driver struct {
}

func (driver *Driver) Open(option driver.Option) (driver.Crawler, error) {
	return &Crawler{Option: option}, nil
}

func init() {
	driver.Register(DriverName, &Driver{})
}

type Crawler struct {
	Option driver.Option
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.odaily.news", nil)
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
	for _, script := range dom.FindAllStrict("script") {
		_body := strings.TrimSpace(script.Text())
		idx := strings.Index(_body, "window.__INITIAL_STATE__ = ")
		if idx < 0 {
			continue
		}
		_body = _body[idx:]
		_body = strings.TrimPrefix(_body, "window.__INITIAL_STATE__ = ")
		_body = strings.Trim(_body, ";")
		body := &body{}
		if err := json.Unmarshal([]byte(_body), body); err != nil {
			return nil, err
		}
		for _, topPost := range body.Home.TopPost {
			url := fmt.Sprintf("https://www.odaily.news/post/%d", topPost.ID)
			board.AppendTitleSummaryURL(topPost.Title, topPost.Summary, url)
		}
		return board, nil
	}
	return nil, fmt.Errorf("not found body")
}

type body struct {
	Home struct {
		TopPost []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Summary string `json:"summary"`
		} `json:"topPost"`
	} `json:"home"`
}
