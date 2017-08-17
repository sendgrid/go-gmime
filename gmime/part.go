package gmime

/*
#cgo pkg-config: gmime-3.0
#include "gmime.h"

*/
import "C"

import (
	"net/textproto"
	"unsafe"
)

// Part is a wrapper for message parts
type Part struct {
	gmimePart *C.GMimeObject
}

// ContentType returns part's content type
func (p *Part) ContentType() string {
	ctype := C.gmime_get_content_type_string(p.gmimePart)
	defer C.g_free(C.gpointer(unsafe.Pointer(ctype)))
	return C.GoString(ctype)
}

// IsText returns true if part's mime is text/*
func (p *Part) IsText() bool {
	return gobool(C.gmime_is_text_part(p.gmimePart))
}

// Text returns text portion of the part if it's mime is text/*
func (p *Part) Text() string {
	content := C.gmime_get_content_string(p.gmimePart)
	defer C.g_free(C.gpointer(unsafe.Pointer(content)))
	return C.GoString(content)
}

// Bytes returns decoded raw bytes of the part, most useful to access attachment data
func (p *Part) Bytes() []byte {
	b := C.gmime_get_bytes(p.gmimePart)
	defer C.g_byte_array_free((*C.GByteArray)(unsafe.Pointer(b)), C.TRUE)
	return C.GoBytes(unsafe.Pointer(b.data), C.int(b.len))
}

// SetText replaces text content if part is text/*
func (p *Part) SetText(text string) error {
	// TODO: Optimize this
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.g_mime_text_part_set_text((*C.GMimeTextPart)(unsafe.Pointer(p.gmimePart)), cstr)
	return nil
}

// SetHeader sets or replaces specified header
func (p *Part) SetHeader(name string, value string) {
	headers := C.g_mime_object_get_header_list(p.asGMimeObject())
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	cCharset := C.CString("UTF-8")
	defer C.free(unsafe.Pointer(cCharset))

	C.g_mime_header_list_set(headers, cName, cValue, cCharset)
}

// Headers gives you all headers for part
func (p *Part) Headers() textproto.MIMEHeader {
	return nil
}

func (p *Part) asGMimeObject() *C.GMimeObject {
	return p.gmimePart
}
