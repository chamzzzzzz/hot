package hot

import (
	"time"
)

type Hot struct {
	Title   string
	Summary string
	Date    time.Time
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

func (b *Board) Append(title, summary string, date time.Time) *Hot {
	hot := &Hot{
		Title:   title,
		Summary: summary,
		Date:    date,
	}
	b.Hots = append(b.Hots, hot)
	return hot
}

type Crawler interface {
	Name() string
	Crawl() (*Board, error)
}

type Archiver interface {
	Name() string
	Archive(*Board) (int, error)
}
