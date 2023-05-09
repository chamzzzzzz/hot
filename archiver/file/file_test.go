package file

import (
	"fmt"
	"testing"
	"time"

	"github.com/chamzzzzzz/hot"
)

func TestArchive(t *testing.T) {
	a := Archiver{}

	board := hot.NewBoard("test")
	for i := 0; i < 1000; i++ {
		board.Append(fmt.Sprintf("Title_%d", i), fmt.Sprintf("Summary_%d", i), time.Now())
		board.Append4(fmt.Sprintf("Title_%d", i), fmt.Sprintf("Summary_%d", i), fmt.Sprintf("URL_%d", i), fmt.Sprintf("Catalog_%d", i))
		board.Append5(fmt.Sprintf("Title_%d", i), fmt.Sprintf("Summary_%d", i), fmt.Sprintf("URL_%d", i), fmt.Sprintf("Catalog_Extra_%d", i), fmt.Sprintf("Extra_%d", i))
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
