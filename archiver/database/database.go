package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chamzzzzzz/hot"
)

type Archiver struct {
	DriverName     string
	DataSourceName string
}

func (a *Archiver) Name() string {
	return "database-archiver"
}

func (a *Archiver) Archive(board *hot.Board) (archived int, err error) {
	if err = a.check(board); err != nil {
		return
	}

	var db *sql.DB
	if db, err = sql.Open(a.DriverName, a.DataSourceName); err != nil {
		return
	}
	defer db.Close()

	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (Date TEXT, PublishDate TEXT, Title TEXT, Summary TEXT, URL TEXT, Catalog TEXT, Extra TEXT);", board.Name)); err != nil {
		return
	}

	var ignore bool
	for _, hot := range board.Hots {
		if ignore, err = a.insert(db, board.Name, hot); err != nil {
			return
		} else if !ignore {
			archived++
		}
	}
	return
}

func (a *Archiver) insert(db *sql.DB, tableName string, hot *hot.Hot) (ignore bool, err error) {
	var rows *sql.Rows
	if rows, err = db.Query(fmt.Sprintf("SELECT Date FROM %s WHERE Title = ? AND Summary = ? AND Catalog = ?", tableName), hot.Title, hot.Summary, hot.Catalog); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		ignore = true
		return
	}

	if _, err = db.Exec(fmt.Sprintf("INSERT INTO %s(Date, PublishDate, Title, Summary, URL, Catalog, Extra) VALUES(?,?,?,?,?,?,?)", tableName), hot.Date.Format(time.RFC3339), hot.PublishDate.Format(time.RFC3339), hot.Title, hot.Summary, hot.URL, hot.Catalog, hot.Extra); err != nil {
		return
	}
	return
}

func (a *Archiver) check(board *hot.Board) error {
	if len(board.Hots) == 0 {
		return fmt.Errorf("empty")
	}

	for _, hot := range board.Hots {
		if hot.URL != "" {
			if !strings.HasPrefix(hot.URL, "http://") && !strings.HasPrefix(hot.URL, "https://") {
				return fmt.Errorf("url `%s` imperfect", hot.URL)
			}
		}
	}
	return nil
}
