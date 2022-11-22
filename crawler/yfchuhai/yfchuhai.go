package yfchuhai

import (
	"github.com/anaskhan96/soup"
	"github.com/chamzzzzzz/hot"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Crawler struct {
}

func (c *Crawler) Name() string {
	return "yfchuhai"
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.yfchuhai.com", nil)
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
	for _, a := range dom.FindAllStrict("a", "class", "list__item flex") {
		h2 := a.Find("h2")
		if h2.Error != nil {
			return nil, h2.Error
		}
		span := a.Find("span", "class", "text")
		if span.Error != nil {
			return nil, span.Error
		}
		title := strings.TrimSpace(h2.Text())
		url := "https://www.yfchuhai.com" + strings.TrimSpace(a.Attrs()["href"])
		date, err := strtotime(strings.TrimSpace(span.Text()))
		if err != nil {
			return nil, err
		}
		board.AppendTitleURLDate(title, url, date)
	}
	return board, nil
}

func strtotime(str string) (time.Time, error) {
	if date, err := time.ParseInLocation("01-02 15:04", str, time.Local); err == nil {
		date = date.AddDate(time.Now().Year(), 0, 0)
		return date, nil
	} else {
		if strings.Contains(str, "小时前") {
			str = strings.TrimSpace(strings.Trim(str, "小时前"))
			if hour, err := strconv.Atoi(str); err != nil {
				return date, err
			} else {
				date = time.Now().Add(time.Hour * time.Duration(-hour))
				return date, nil
			}
		} else if strings.Contains(str, "天前") {
			str = strings.TrimSpace(strings.Trim(str, "天前"))
			if day, err := strconv.Atoi(str); err != nil {
				return date, err
			} else {
				date = time.Now().Add(time.Hour * 24 * time.Duration(-day))
				return date, nil
			}
		} else if strings.Contains(str, "前天") {
			str = strings.TrimSpace(strings.Trim(str, "前天"))
			if date, err := time.ParseInLocation("15:04", str, time.Local); err != nil {
				return date, err
			} else {
				date = date.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day()-1)
				date = date.Add(time.Hour * 24 * time.Duration(-2))
				return date, nil
			}
		} else if strings.Contains(str, "昨天") {
			str = strings.TrimSpace(strings.Trim(str, "昨天"))
			if date, err := time.ParseInLocation("15:04", str, time.Local); err != nil {
				return date, err
			} else {
				date = date.AddDate(time.Now().Year(), int(time.Now().Month())-1, time.Now().Day()-1)
				date = date.Add(time.Hour * 24 * time.Duration(-1))
				return date, nil
			}
		} else {
			return date, err
		}
	}
}
