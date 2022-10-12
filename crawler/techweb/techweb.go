package techweb

import (
	"encoding/xml"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"time"
)

type Crawler struct {
	Proxy string
}

func (c *Crawler) Name() string {
	return "techweb"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://www.techweb.com.cn/rss/hotnews.xml", nil)
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
		url := item.Link
		date, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			return nil, err
		}
		if title == "" {
			continue
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}

type body struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Item []struct {
			Title   string `xml:"title"`
			Link    string `xml:"link"`
			PubDate string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}
