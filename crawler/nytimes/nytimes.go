package nytimes

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Crawler struct {
	Proxy string
}

func (c *Crawler) Name() string {
	return "nytimes"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	if c.Proxy != "" {
		if proxyUrl, err := url.Parse(c.Proxy); err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}

	req, err := http.NewRequest("GET", "https://cn.nytimes.com/async/mostviewed/all/?lang=zh-hans", nil)
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

	body := &body{}
	if err := json.Unmarshal(data, body); err != nil {
		return nil, err
	} else if body.Code != 0 {
		return nil, fmt.Errorf("body code: %d", body.Code)
	}

	board := hot.NewBoard(c.Name())
	for _, daily := range body.List.Daily {
		board.AppendTitleSummaryURL(daily.Headline, daily.Summary, daily.URL)
	}
	return board, nil
}

type body struct {
	Code int `json:"code"`
	List struct {
		Daily []struct {
			URL           string `json:"url"`
			Headline      string `json:"headline"`
			ShortHeadline string `json:"short_headline"`
			Summary       string `json:"summary"`
		} `json:"daily"`
		Weekly []struct {
			URL           string `json:"url"`
			Headline      string `json:"headline"`
			ShortHeadline string `json:"short_headline"`
			Summary       string `json:"summary"`
		} `json:"weekly"`
		WeeklySlideshow []struct {
			URL           string `json:"url"`
			Headline      string `json:"headline"`
			ShortHeadline string `json:"short_headline"`
			Summary       string `json:"summary"`
		} `json:"weekly_slideshow"`
	} `json:"list"`
}
