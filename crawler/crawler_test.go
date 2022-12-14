package crawler

import (
	"testing"
)

func TestDrivers(t *testing.T) {
	drivers := Drivers()
	t.Log("driver count:", len(drivers))
	t.Log(drivers)
}
