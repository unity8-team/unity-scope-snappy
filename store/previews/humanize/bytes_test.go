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
