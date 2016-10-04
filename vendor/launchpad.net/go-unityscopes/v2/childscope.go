package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"runtime"
	"unsafe"
)

type ChildScope struct {
	c *C._ChildScope
}

func finalizeChildScope(childscope *ChildScope) {
	if childscope.c != nil {
		C.destroy_child_scope((*C._ChildScope)(childscope.c))
	}
	childscope.c = nil
}

func makeChildScope(c *C._ChildScope) *ChildScope {
	childscope := new(ChildScope)
	runtime.SetFinalizer(childscope, finalizeChildScope)
	childscope.c = (*C._ChildScope)(c)
	return childscope
}

// NewChildScope creates a new ChildScope with the given id, metadata, enabled state and keywords
func NewChildScope(id string, metadata *ScopeMetadata, enabled bool, keywords []string) *ChildScope {
	var cEnabled C.int
	if enabled {
		cEnabled = 1
	} else {
		cEnabled = 0
	}
	return makeChildScope(C.new_child_scope(strData(id),
		(*C._ScopeMetadata)(metadata.m),
		cEnabled,
		joinedStrData(keywords)))
}

// Id returns the identifier of the child scope
func (childscope *ChildScope) Id() string {
	s := C.child_scope_get_id(childscope.c)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}
