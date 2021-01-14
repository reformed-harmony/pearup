package algorithm

import (
	"testing"
)

func TestCopyMap(t *testing.T) {
	var (
		a = map[int64]*algUser{0: nil}
		b = copyMap(a)
	)
	delete(b, 0)
	if _, ok := a[0]; !ok {
		t.Fatalf("key missing from map")
	}
}
