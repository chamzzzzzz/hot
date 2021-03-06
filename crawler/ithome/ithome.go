package ithome

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BoardGame = "game"
)

type Crawler struct {
	BoardName string
}

func (c *Crawler) Name() string {
	if c.BoardName == BoardGame {
		return "ithome_x_game"
	}
	return "ithome"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.ithome.com/block/rank.html?d=game", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	html, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(html))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	date := time.Now()
	ulId := "d-1"
	if c.BoardName == BoardGame {
		ulId = "d-4"
	}
	ul := dom.FindStrict("ul", "id", ulId)
	if ul.Error != nil {
		return nil, ul.Error
	}
	for _, a := range ul.FindAllStrict("a") {
		title := a.Text()
		board.Append(title, "", date)
	}
	return board, nil
}
