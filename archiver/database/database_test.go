package database

import (
	"fmt"
	"github.com/chamzzzzzz/hot"
	"os"
	"testing"
	"time"
)

func TestArchive(t *testing.T) {
	a := Archiver{
		DriverName:     os.Getenv("HOT_ARCHIVER_DATABASE_TEST_DRIVER_NAME"),
		DataSourceName: os.Getenv("HOT_ARCHIVER_DATABASE_TEST_DATA_SOURCE_NAME"),
	}

	board := hot.NewBoard("test")
	for i := 0; i < 1000; i++ {
		board.Append(fmt.Sprintf("Title_%d", i), fmt.Sprintf("Content_%d", i), time.Now())
	}

	if archived, err := a.Archive(board); err != nil {
		t.Error(err)
	} else {
		t.Logf("archived=%d\n", archived)
	}
}
