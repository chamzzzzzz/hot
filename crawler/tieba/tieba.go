package tieba

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "tieba"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://tieba.baidu.com/hottopic/browse/topicList", nil)
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
	} else if body.Errno != 0 {
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
