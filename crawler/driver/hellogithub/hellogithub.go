package hellogithub

import (
	"fmt"
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
	"github.com/chamzzzzzz/hot/crawler/driver"
	"github.com/chamzzzzzz/hot/crawler/httputil"
)

const (
	DriverName  = "hellogithub"
	ProxySwitch = false
	URL         = "https://api.hellogithub.com/v1/?sort_by=hot"
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
	body := &body{}
	if err := httputil.Request("GET", URL, nil, "json", body, httputil.NewOption(c.Option, ProxySwitch)); err != nil {
		return nil, err
	}

	board := hot.NewBoard(c.Name())
	for _, data := range body.Data {
		t, err := time.ParseInLocation("2006-01-02T15:04:05", data.UpdatedAt, time.Local)
		if err != nil {
			return nil, err
		}
		hot := &hot.Hot{
			Title:       strings.TrimSpace(data.Title),
			Summary:     strings.TrimSpace(data.Summary),
			URL:         fmt.Sprintf("https://hellogithub.com/repository/%s", data.ItemID),
			PublishDate: t,
		}
		board.Append(hot)
	}
	return board, nil
}

type body struct {
	Success bool `json:"success"`
	Data    []struct {
		ItemID    string `json:"item_id"`
		Title     string `json:"title"`
		Summary   string `json:"summary"`
		IsHot     bool   `json:"is_hot"`
		UpdatedAt string `json:"updated_at"`
	} `json:"data"`
}

func (body *body) NormalizedCode() int {
	if body.Success {
		return 0
	} else {
		return 1
	}
}
