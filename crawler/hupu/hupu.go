package hupu

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	Basketball = "basketball"
	Football   = "football"
	Gambia     = "gambia"
)

type Crawler struct {
	BoardName string
}

func (c *Crawler) Name() string {
	switch c.BoardName {
	case Basketball:
		return "hupu_x_basketball"
	case Football:
		return "hupu_x_football"
	case Gambia:
		return "hupu_x_gambia"
	default:
		return "hupu"
	}
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.BoardName {
	case Basketball:
		return c.hupu_x_basketball()
	case Football:
		return c.hupu_x_football()
	case Gambia:
		return c.hupu_x_gambia()
	default:
		return c.hupu()
	}
}

func (c *Crawler) hupu_x_basketball() (*hot.Board, error) {
	return c.hupuWithClasses([]string{"newpcbasketball"})
}

func (c *Crawler) hupu_x_football() (*hot.Board, error) {
	return c.hupuWithClasses([]string{"newpcsoccer"})
}

func (c *Crawler) hupu_x_gambia() (*hot.Board, error) {
	return c.hupuWithClasses([]string{"newpcbbs"})
}

func (c *Crawler) hupu() (*hot.Board, error) {
	return c.hupuWithClasses([]string{"newpcbasketball", "newpcsoccer", "newpcbbs"})
}

func (c *Crawler) hupuWithClasses(hotClasses []string) (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.hupu.com", nil)
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

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	for _, hotClass := range hotClasses {
		for _, a := range dom.FindAllStrict("a", "class", hotClass) {
			div := a.FindStrict("div", "class", "hot-title")
			if div.Error != nil {
				return nil, div.Error
			}
			title := strings.TrimSpace(div.Text())
			url := strings.TrimSpace(a.Attrs()["href"])
			board.AppendTitleURL(title, url)
		}
	}
	return board, nil
}
