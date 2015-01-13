package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"
import (
	"unsafe"
)

type ContentEncoding interface {
	ToString() string
}

type rawContentEncoding interface {
	ContentEncoding
	rawContentEncoding() C.GMimeContentEncoding
}

type aContentEncoding struct {
	encoding C.GMimeContentEncoding
}

func CastContentEncoding(encoding C.GMimeContentEncoding) ContentEncoding {
	return &aContentEncoding{encoding: encoding}
}

func NewContentEncodingFromString(str string) ContentEncoding {
	var _str *C.char = C.CString(str)
	defer C.free(unsafe.Pointer(_str))
	encoding := C.g_mime_content_encoding_from_string(_str)
	return CastContentEncoding(encoding)
}

func (e *aContentEncoding) ToString() string {
	var _str *C.char = C.g_mime_content_encoding_to_string(e.encoding)
	return C.GoString(_str)
}

func (e *aContentEncoding) rawContentEncoding() C.GMimeContentEncoding {
	return e.encoding
}
