package store

import (
	"testing"
)

// Data for TestHumanizeBytes
var humanizeBytesTests = []struct {
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

func TestHumanizeBytes(t *testing.T) {
	for i, test := range humanizeBytesTests {
		humanized := humanizeBytes(test.size)
		if humanized != test.expectedHumanized {
			t.Errorf("Test case %d: Got %s, expected %s", i, humanized, test.expectedHumanized)
		}
	}
}
