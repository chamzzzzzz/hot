package wsj

import (
	"encoding/xml"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Crawler struct {
	Proxy string
}

func (c *Crawler) Name() string {
	return "wsj"
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

	req, err := http.NewRequest("GET", "https://cn.wsj.com/zh-hans/rss", nil)
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
	if err := xml.Unmarshal(data, body); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, item := range body.Channel.Item {
		title := item.Title
		summary := item.Description
		url := item.Link
		date, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return nil, err
		}
		date = time.Unix(date.Unix(), 0)
		i := strings.Index(summary, "<p>")
		j := strings.Index(summary, "</p>")
		if i >= 0 && j >= 0 {
			summary = summary[i+3 : j]
		}
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Channel struct {
		Item []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}
