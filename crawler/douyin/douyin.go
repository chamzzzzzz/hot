package douyin

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
	return "douyin"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://aweme.snssdk.com/aweme/v1/hot/search/list/", nil)
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
	} else if bodyJson.StatusCode != 0 {
		return nil, fmt.Errorf("body status_code: %d", bodyJson.StatusCode)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, word := range bodyJson.Data.WordList {
		board.Append(word.Word, "", date)
	}
	return board, nil
}

type bodyJson struct {
	StatusCode int `json:"status_code"`
	Data       struct {
		WordList []struct {
			Word string `json:"word"`
		} `json:"word_list"`
	} `json:"data"`
}
