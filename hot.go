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

func (b *Board) Append0(hot *Hot) *Hot {
	b.Hots = append(b.Hots, hot)
	return hot
}

func (b *Board) Append(title, summary string, date time.Time) *Hot {
	return b.Append5x2(title, summary, "", "", "", date, time.Now())
}

func (b *Board) Append1(title string) *Hot {
	return b.Append5x2(title, "", "", "", "", time.Now(), time.Now())
}

func (b *Board) Append1x1(title string, publishDate time.Time) *Hot {
	return b.Append5x2(title, "", "", "", "", time.Now(), publishDate)
}

func (b *Board) Append2(title, summary string) *Hot {
	return b.Append5x2(title, summary, "", "", "", time.Now(), time.Now())
}

func (b *Board) Append2x1(title, summary string, publishDate time.Time) *Hot {
	return b.Append5x2(title, summary, "", "", "", time.Now(), publishDate)
}

func (b *Board) Append3(title, summary, url string) *Hot {
	return b.Append5x2(title, summary, url, "", "", time.Now(), time.Now())
}

func (b *Board) Append3x1(title, summary, url string, publishDate time.Time) *Hot {
	return b.Append5x2(title, summary, url, "", "", time.Now(), publishDate)
}

func (b *Board) Append4(title, summary, url, catalog string) *Hot {
	return b.Append5x2(title, summary, url, catalog, "", time.Now(), time.Now())
}

func (b *Board) Append4x1(title, summary, url, catalog string, publishDate time.Time) *Hot {
	return b.Append5x2(title, summary, url, catalog, "", time.Now(), publishDate)
}

func (b *Board) Append5(title, summary, url, catalog, extra string) *Hot {
	return b.Append5x2(title, summary, url, catalog, extra, time.Now(), time.Now())
}

func (b *Board) Append5x1(title, summary, url, catalog, extra string, publishDate time.Time) *Hot {
	return b.Append5x2(title, summary, url, catalog, extra, time.Now(), publishDate)
}

func (b *Board) AppendTitleSummaryURL(title, summary, url string) *Hot {
	return b.Append3(title, summary, url)
}

func (b *Board) AppendTitleURL(title, url string) *Hot {
	return b.Append3(title, "", url)
}

func (b *Board) AppendTitleURLDate(title, url string, publishDate time.Time) *Hot {
	return b.Append3x1(title, "", url, publishDate)
}

func (b *Board) Append5x2(title, summary, url, catalog, extra string, date, publishDate time.Time) *Hot {
	hot := &Hot{
		Title:       title,
		Summary:     summary,
		URL:         url,
		Catalog:     catalog,
		Extra:       extra,
		Date:        date,
		PublishDate: publishDate,
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
