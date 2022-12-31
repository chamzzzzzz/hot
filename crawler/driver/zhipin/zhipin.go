package zhipin

import (
	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
	"strings"
)

const (
	Question         = "question"
	HotSearchJob     = "hotsearchjob"
	HotRecruitingJob = "hotrecruitingjob"
)

const (
	DriverName          = "zhipin"
	ProxySwitch         = false
	URL                 = "https://youle.zhipin.com/recommend/selected/"
	QuestionURL         = URL
	HotSearchJobURL     = "https://baike.zhipin.com/wiki/"
	HotRecruitingJobURL = "https://www.zhipin.com/"
)

type Driver struct {
}

func (driver *Driver) Open(option driver.Option) (driver.Crawler, error) {
	return &Crawler{Option: option}, nil
}

func init() {
	driver.Register(DriverName, &Driver{})
}

type Crawler struct {
	Option driver.Option
}

func (c *Crawler) Driver() driver.Driver {
	return &Driver{}
}

func (c *Crawler) Name() string {
	return DriverName
}

func (c *Crawler) Crawl() (*hot.Board, error) {
	switch c.Option.Catalog {
	case Question:
		return c.question(nil)
	case HotSearchJob:
		return c.hotsearchjob(nil)
	case HotRecruitingJob:
		return c.hotrecruitingjob(nil)
	default:
		return c.all()
	}
}

func (c *Crawler) all() (*hot.Board, error) {
	board := hot.NewBoard(c.Name())
	if _, err := c.question(board); err != nil {
		return nil, err
	}
	if _, err := c.hotsearchjob(board); err != nil {
		return nil, err
	}
	if _, err := c.hotrecruitingjob(board); err != nil {
		return nil, err
	}
	return board, nil
}

func (c *Crawler) question(board *hot.Board) (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", QuestionURL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	for _, li := range dom.QueryAll("li", "class", "hot-item") {
		a, err := li.Find("a")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Text())
		url := strings.TrimSpace(a.Href())
		board.Append4(title, "", url, Question)
	}
	return board, nil
}

func (c *Crawler) hotsearchjob(board *hot.Board) (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", HotSearchJobURL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	div, err := dom.Find("div", "class", "baike-hot-list-item hotList")
	if err != nil {
		return nil, err
	}
	for _, li := range div.QueryAll("li", "class", "list-item") {
		a, err := li.Find("a")
		if err != nil {
			return nil, err
		}
		span, err := li.Find("span")
		if err != nil {
			return nil, err
		}
		title := strings.TrimSpace(a.Text())
		summary := strings.TrimSpace(span.Text())
		url := strings.TrimSpace(a.Href())
		board.Append4(title, summary, url, HotSearchJob)
	}
	return board, nil
}

func (c *Crawler) hotrecruitingjob(board *hot.Board) (*hot.Board, error) {
	dom := &httputil.DOM{}
	if err := httputil.Request("GET", HotRecruitingJobURL, nil, "dom", dom, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	if board == nil {
		board = hot.NewBoard(c.Name())
	}
	div, err := dom.Find("div", "class", "common-tab-box merge-city-job hot-job-box")
	if err != nil {
		return nil, err
	}
	for _, li := range div.QueryAll("li") {
		a, err := li.Find("a", "class", "job-info")
		if err != nil {
			return nil, err
		}
		a2, err := li.Find("a", "class", "user-info")
		if err != nil {
			return nil, err
		}
		title := strings.Join(strings.Fields(strings.TrimSpace(a.FullText())), "---")
		summary := strings.Join(strings.Fields(strings.TrimSpace(a2.FullText())), "---")
		url := "https://www.zhipin.com" + strings.TrimSpace(a.Href())
		board.Append4(title, summary, url, HotRecruitingJob)
	}
	return board, nil
}
