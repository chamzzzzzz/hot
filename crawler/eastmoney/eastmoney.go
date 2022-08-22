package eastmoney

import (
	"encoding/json"
	"fmt"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var re = regexp.MustCompile(`//searchapi\.eastmoney\.com/api/hotkeyword/get\?count=20&token=([A-Z0-9]+)`)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "eastmoney"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	URL, err := c.getURL()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
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
	} else if bodyJson.Status != 0 {
		return nil, fmt.Errorf("body status: %d", bodyJson.Status)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson.Data {
		board.Append(data.KeyPhrase, "", date)
	}
	return board, nil
}

func (c *Crawler) getURL() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://so.eastmoney.com/newstatic/js/page/welcome.js", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	URL := re.FindString(string(body))
	if URL == "" {
		return "", fmt.Errorf("not match")
	}
	return "https:" + URL, nil
}

type bodyJson struct {
	Data []struct {
		KeyPhrase string `json:"KeyPhrase"`
	} `json:"Data"`
	Status  int    `json:"Status"`
	Message string `json:"Message"`
}
