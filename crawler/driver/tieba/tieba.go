package tieba

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
	"time"
)

const (
	DriverName  = "tieba"
	ProxySwitch = false
	URL         = "http://tieba.baidu.com/hottopic/browse/topicList"
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
	body := &body{}
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}
	if body.Errno != 0 {
		return nil, fmt.Errorf("body errno: %d", body.Errno)
	}

	board := hot.NewBoard(c.Name())
	for _, topic := range body.Data.BangTopic.TopicList {
		title := topic.TopicName
		summary := topic.TopicDesc
		url := strings.ReplaceAll(topic.TopicURL, "amp;", "")
		date := time.Unix(topic.CreateTime, 0)
		board.Append3x1(title, summary, url, date)
	}
	return board, nil
}

type body struct {
	Data struct {
		BangTopic struct {
			TopicList []struct {
				TopicName  string `json:"topic_name"`
				TopicDesc  string `json:"topic_desc"`
				TopicURL   string `json:"topic_url"`
				CreateTime int64  `json:"create_time"`
			} `json:"topic_list"`
		} `json:"bang_topic"`
	} `json:"data"`
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}
