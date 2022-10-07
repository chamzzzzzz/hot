package jinse

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
	return "jinse"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	return c.jinse()
}

func (c *Crawler) jinse_x_search() (*hot.Board, error) {
	return c.jinseWithURL("https://newapi.jinse.com/noah/v1/hot-search")
}

func (c *Crawler) jinse_x_article() (*hot.Board, error) {
	return c.jinseWithURL("https://newapi.jinse.com/noah/v1/articles/hot?hour_diff=24")
}

func (c *Crawler) jinse() (*hot.Board, error) {
	board, err := c.jinse_x_search()
	if err != nil {
		return nil, err
	}

	board2, err := c.jinse_x_article()
	if err != nil {
		return nil, err
	}

	for _, hot := range board2.Hots {
		board.Append(hot.Title, hot.Summary, hot.Date)
	}
	return board, nil
}

func (c *Crawler) jinseWithURL(URL string) (*hot.Board, error) {
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
	} else if bodyJson.StatusCode != 0 {
		return nil, fmt.Errorf("body status_code: %d", bodyJson.StatusCode)
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	for _, data := range bodyJson.Data {
		board.Append(data.Title, data.JumpURL, date)
	}
	return board, nil
}

type bodyJson struct {
	StatusCode int `json:"status_code"`
	Data       []struct {
		Title   string `json:"title"`
		JumpURL string `json:"jump_url"`
	} `json:"data"`
}
