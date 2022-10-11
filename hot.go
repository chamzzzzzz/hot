package hot

import (
	"time"
)

type Hot struct {
	Title    string
	Summary  string
	URL      string
	Catagory string
	Date     time.Time
}

type Board struct {
	Name string
	Hots []*Hot
}

func NewBoard(name string) *Board {
	return &Board{
		Name: name,
	}
}

func (b *Board) Append0(hot *Hot) *Hot {
	b.Hots = append(b.Hots, hot)
	return hot
}

func (b *Board) Append(title, summary string, date time.Time) *Hot {
	return b.Append3(title, summary, date)
}

func (b *Board) Append5(title, summary, url, catagory string, date time.Time) *Hot {
	hot := &Hot{
		Title:    title,
		Summary:  summary,
		URL:      url,
		Catagory: catagory,
		Date:     date,
	}
	return b.Append0(hot)
}

func (b *Board) Append4(title, summary, url string, date time.Time) *Hot {
	hot := &Hot{
		Title:   title,
		Summary: summary,
		URL:     url,
		Date:    date,
	}
	return b.Append0(hot)
}

func (b *Board) Append3(title, summary string, date time.Time) *Hot {
	hot := &Hot{
		Title:   title,
		Summary: summary,
		Date:    date,
	}
	return b.Append0(hot)
}

func (b *Board) Append2(title string, date time.Time) *Hot {
	hot := &Hot{
		Title: title,
		Date:  date,
	}
	return b.Append0(hot)
}

type Crawler interface {
	Name() string
	Crawl() (*Board, error)
}

type Archiver interface {
	Name() string
	Archive(*Board) (int, error)
}
