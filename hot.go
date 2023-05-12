package hot

import (
	"fmt"
	"time"
)

type Hot struct {
	Title       string
	Summary     string
	URL         string
	Catalog     string
	Extra       string
	Date        time.Time
	PublishDate time.Time
}

func (hot *Hot) String() string {
	return fmt.Sprintf("%s | %s | %s | %s | %s | %s | %s", hot.Title, hot.Summary, hot.URL, hot.Catalog, hot.Extra, hot.Date.Format(time.RFC3339), hot.PublishDate.Format(time.RFC3339))
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

func (b *Board) Append(hot *Hot) *Hot {
	if hot.Date.IsZero() {
		hot.Date = time.Now()
	}
	if hot.PublishDate.IsZero() {
		hot.PublishDate = time.Now()
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
