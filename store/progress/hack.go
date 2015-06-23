package progress

import "launchpad.net/unity-scope-snappy/webdm"

// Hack is a workaround for having no concept of progress in this scope. Until a
// decent method has been devised, this struct holds the information necessary
// to display a placeholder widget for manual refreshing.
type Hack struct {
	DesiredStatus webdm.Status
}
