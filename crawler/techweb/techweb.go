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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var bodyXML bodyXML
	if err := xml.Unmarshal(body, &bodyXML); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, item := range bodyXML.Channel.Item {
		title := item.Title
		summary := item.Link
		if title == "" {
			continue
		}
		board.Append(title, summary, date)
	}
	return board, nil
}

type bodyXML struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Item []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}
