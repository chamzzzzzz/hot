package infoq

import (
	"bytes"
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
	return "infoq"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://www.infoq.cn/public/v1/article/getHotList", bytes.NewReader([]byte(`{"score":null,"type":1,"size":30}`)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.Set("Origin", "https://www.infoq.cn")
	req.Header.Set("Content-Type", "application/json")

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
	} else if bodyJson.Code != 0 {
		return nil, fmt.Errorf("body code: %d", bodyJson.Code)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson.Data {
		board.Append(data.ArticleTitle, data.ArticleSummary, date)
	}
	return board, nil
}

type bodyJson struct {
	Code int `json:"code"`
	Data []struct {
		ArticleSummary string `json:"article_summary"`
		ArticleTitle   string `json:"article_title"`
	} `json:"data"`
}
