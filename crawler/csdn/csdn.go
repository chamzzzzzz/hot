package csdn

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
	return "csdn"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.csdn.net", nil)
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
		if strings.HasPrefix(body, "window.__INITIAL_STATE__= ") {
			body = strings.TrimPrefix(body, "window.__INITIAL_STATE__= ")
			body = strings.Trim(body, ";")
			bodyJson := &bodyJson{}
			if err := json.Unmarshal([]byte(body), bodyJson); err != nil {
				return nil, err
			}
			for _, headline := range bodyJson.PageData.Data.WwwHeadlines {
				board.Append(fmt.Sprintf("%s | %s", headline.Title, headline.Description), headline.URL, date)
			}
			for _, headhot := range bodyJson.PageData.Data.WwwHeadhot {
				board.Append(fmt.Sprintf("%s | %s", headhot.Title, headhot.Description), headhot.URL, date)
			}
			return board, nil
		}
	}
	return nil, fmt.Errorf("not found body")
}

type bodyJson struct {
	PageData struct {
		Data struct {
			WwwHeadlines []struct {
				Description string `json:"description"`
				Title       string `json:"title"`
				URL         string `json:"url"`
			} `json:"www-Headlines"`
			WwwHeadhot []struct {
				Description string `json:"description"`
				Title       string `json:"title"`
				URL         string `json:"url"`
			} `json:"www-headhot"`
		} `json:"data"`
	} `json:"pageData"`
}
