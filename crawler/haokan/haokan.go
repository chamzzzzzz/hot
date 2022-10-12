package haokan

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "haokan"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://haokan.baidu.com/videoui/api/hotwords", nil)
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
	for _, hotword := range body.Data.Response.Hotwords {
		board.Append1(hotword)
	}
	return board, nil
}

type body struct {
	Data struct {
		Response struct {
			Hotwords []string `json:"hotwords"`
		} `json:"response"`
	} `json:"data"`
	Errno int    `json:"errno"`
	Error string `json:"error"`
}
