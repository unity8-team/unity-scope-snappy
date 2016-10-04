/* Copyright (C) 2015 Canonical Ltd.
 *
 * This file is part of unity-scope-snappy.
 *
 * unity-scope-snappy is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * unity-scope-snappy is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * unity-scope-snappy. If not, see <http://www.gnu.org/licenses/>.
 */

package humanize

import (
	"github.com/dustin/go-humanize"
	"regexp"
)

var regex = regexp.MustCompile("([^a-zA-Z]*)([a-zA-Z]*)")

// Bytes translates a raw number of bytes into a nice human-readable
// representation. This function was written because go-humanize doesn't put
// spaces between the number and the unit.
//
// NOTE: This function uses SI units (e.g. 1 kB == 1000 bytes)
//
// Parameters:
// bytes: Number of bytes
//
// Returns:
// - Humanized number in a string ("Unknown" if `bytes` is negative)
func Bytes(bytes int64) string {
	if bytes < 0 {
		return "Unknown"
	}

	return regex.ReplaceAllString(humanize.Bytes(uint64(bytes)), "$1 $2")
}
