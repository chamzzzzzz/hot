package database

import (
	"database/sql"
	"fmt"
	"github.com/chamzzzzzz/hot"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Archiver struct {
	DriverName     string
	DataSourceName string
}

func (a *Archiver) Name() string {
	return "database"
}

func (a *Archiver) Archive(board *hot.Board) (archived int, err error) {
	var db *sql.DB
	if db, err = sql.Open(a.DriverName, a.DataSourceName); err != nil {
		return
	}
	defer db.Close()

	if _, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (Date TEXT, Title TEXT, Summary TEXT);", board.Name)); err != nil {
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
	if rows, err = db.Query(fmt.Sprintf("SELECT Date, Title FROM %s WHERE Title = ?", tableName), hot.Title); err != nil {
		return
	}
	defer rows.Close()

	var date, title string
	var firstTime time.Time
	for rows.Next() {
		if err = rows.Scan(&date, &title); err != nil {
			return
		}

		if firstTime, err = time.Parse(time.RFC3339, date); err != nil {
			return
		}

		if hot.Date.Sub(firstTime) < 24*7*time.Hour {
			ignore = true
			return
		}
	}

	if _, err = db.Exec(fmt.Sprintf("INSERT INTO %s(Date, Title, Summary) VALUES(?,?,?)", tableName), hot.Date.Format(time.RFC3339), hot.Title, hot.Summary); err != nil {
		return
	}
	return
}
