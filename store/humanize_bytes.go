package store

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/dustin/go-humanize"
	"regexp"
)

var regex = regexp.MustCompile("([^a-zA-Z]*)([a-zA-Z]*)")

// humanizeBytes translates a raw number of bytes into a nice human-readable
// representation. This function was written because go-humanize doesn't put
// spaces between the number and the unit.
//
// NOTE: This function uses SI units (e.g. 1 KB == 1000 bytes)
//
// Parameters:
// bytes: Number of bytes
//
// Returns:
// - Humanized number in a string ("Unknown" if `bytes` is negative)
func humanizeBytes(bytes int64) string {
	if bytes < 0 {
		return "Unknown"
	}

	return regex.ReplaceAllString(humanize.Bytes(uint64(bytes)), "$1 $2")
}
