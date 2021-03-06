package tieba

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyJson := &bodyJson{}
	if err := json.Unmarshal(body, bodyJson); err != nil {
		return nil, err
	} else if bodyJson.Errno != 0 {
		return nil, fmt.Errorf("body errno: %d", bodyJson.Errno)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, topic := range bodyJson.Data.BangTopic.TopicList {
		board.Append(topic.TopicName, topic.TopicDesc, date)
	}
	return board, nil
}

type bodyJson struct {
	Data struct {
		BangTopic struct {
			TopicList []struct {
				TopicName string `json:"topic_name"`
				TopicDesc string `json:"topic_desc"`
			} `json:"topic_list"`
		} `json:"bang_topic"`
	} `json:"data"`
	Errno  int    `json:"errno"`
	Errmsg string `json:"errmsg"`
}
