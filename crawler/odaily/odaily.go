package odaily

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "odaily"
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
	for _, script := range dom.FindAllStrict("script") {
		body := strings.TrimSpace(script.Text())
		idx := strings.Index(body, "window.__INITIAL_STATE__ = ")
		if idx < 0 {
			continue
		}
		body = body[idx:]
		body = strings.TrimPrefix(body, "window.__INITIAL_STATE__ = ")
		body = strings.Trim(body, ";")
		bodyJson := &bodyJson{}
		if err := json.Unmarshal([]byte(body), bodyJson); err != nil {
			return nil, err
		}
		for _, topPost := range bodyJson.Home.TopPost {
			board.Append(fmt.Sprintf("%s | %s", topPost.Title, topPost.Summary), fmt.Sprintf("/post/%d", topPost.ID), date)
		}
		return board, nil
	}
	return nil, fmt.Errorf("not found body")
}

type bodyJson struct {
	Home struct {
		TopPost []struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Summary string `json:"summary"`
		} `json:"topPost"`
	} `json:"home"`
}
