package humanize

import (
	"testing"
)

// Data for TestBytes
var bytesTests = []struct {
	size              int64
	expectedHumanized string
}{
	{-1, "Unknown"},
	{0, "0 B"},
	{1024, "1.0 kB"},
	{1485, "1.5 kB"},
	{9999, "10 kB"},
	{1520435, "1.5 MB"},
	{1556925645, "1.6 GB"},
}

func TestBytes(t *testing.T) {
	for i, test := range bytesTests {
		humanized := Bytes(test.size)
		if humanized != test.expectedHumanized {
			t.Errorf("Test case %d: Got %s, expected %s", i, humanized, test.expectedHumanized)
		}
	}
}
