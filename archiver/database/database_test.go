package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/chamzzzzzz/hot"
)

func TestArchive(t *testing.T) {
	a := Archiver{
		DriverName:     os.Getenv("HOT_ARCHIVER_DATABASE_TEST_DRIVER_NAME"),
		DataSourceName: os.Getenv("HOT_ARCHIVER_DATABASE_TEST_DATA_SOURCE_NAME"),
	}

	board := hot.NewBoard("test")
	for i := 0; i < 1000; i++ {
		h := &hot.Hot{
			Title:   fmt.Sprintf("Title_%d", i),
			Summary: fmt.Sprintf("Summary_%d", i),
		}
		board.Append(h)

		h = &hot.Hot{
			Title:   fmt.Sprintf("Title_%d", i),
			Summary: fmt.Sprintf("Summary_%d", i),
			URL:     fmt.Sprintf("URL_%d", i),
			Catalog: fmt.Sprintf("Catalog_%d", i),
		}
		board.Append(h)

		h = &hot.Hot{
			Title:   fmt.Sprintf("Title_%d", i),
			Summary: fmt.Sprintf("Summary_%d", i),
			URL:     fmt.Sprintf("URL_%d", i),
			Catalog: fmt.Sprintf("Catalog_Extra_%d", i),
			Extra:   fmt.Sprintf("Extra_%d", i),
		}
		board.Append(h)
	}

	if archived, err := a.Archive(board); err != nil {
		t.Error(err)
	} else {
		t.Logf("archived=%d\n", archived)
	}

	if archived, err := a.Archive(board); err != nil {
		t.Error(err)
	} else {
		t.Logf("archive again archived=%d\n", archived)
	}
}
