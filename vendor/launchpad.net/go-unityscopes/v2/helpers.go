package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"reflect"
	"unsafe"
)

// strData transforms a string to a form that can be passed to cgo
// without copying data.
func strData(s string) C.StrData {
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return C.StrData{
		data: (*C.char)(unsafe.Pointer(h.Data)),
		length: C.long(len(s)),
	}
}

// byteData transforms a byte array into the same format strData() produces.
func byteData(b []byte) C.StrData {
	if len(b) == 0 {
		return C.StrData{
			data: nil,
			length: 0,
		}
	}
	return C.StrData{
		data: (*C.char)(unsafe.Pointer(&b[0])),
		length: C.long(len(b)),
	}
}


func joinedStrData(a []string) C.StrData {
	total := 0
	for _, s := range a {
		total += len(s) + 1
	}
	buf := make([]byte, total)
	i := 0
	for _, s := range a {
		copy(buf[i:i+len(s)], s)
		buf[i+len(s)] = 0
		i += len(s) + 1
	}
	return byteData(buf)
}
