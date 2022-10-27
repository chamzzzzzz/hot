package daniu

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "daniu"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.daniu523.com/misc.php?mod=ranklist&type=thread&view=heats&orderby=today", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(transform.NewReader(res.Body, simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return nil, err
	}

	dom := soup.HTMLParse(string(data))
	if dom.Error != nil {
		return nil, dom.Error
	}

	board := hot.NewBoard(c.Name())
	div := dom.FindStrict("div", "class", "tl")
	if div.Error != nil {
		return nil, div.Error
	}
	tbody := div.FindAllStrict("tbody")
	if len(tbody) != 2 {
		return nil, fmt.Errorf("tbody count invalid")
	}
	for _, tr := range tbody[1].FindAllStrict("tr") {
		th := tr.Find("th")
		if th.Error != nil {
			return nil, th.Error
		}
		em := tr.FindStrict("em")
		if em.Error != nil {
			fmt.Println(tr.FullText())
			return nil, em.Error
		}
		a := th.Find("a")
		if a.Error != nil {
			return nil, a.Error
		}

		title := strings.TrimSpace(a.Text())
		url := "https://www.daniu523.com/" + strings.TrimSpace(a.Attrs()["href"])
		date, err := time.ParseInLocation("2006-01-02 15:04", strings.TrimSpace(em.Text()), time.Local)
		if err != nil {
			return nil, err
		}
		board.Append3x1(title, "", url, date)
	}
	return board, nil
}
