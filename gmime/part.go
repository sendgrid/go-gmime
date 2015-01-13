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

type Part interface {
	Object
	ContentObject() DataWrapper
	SetContentObject(DataWrapper)
	Filename() string
	Description() string
	ContentLocation() string
	ContentEncoding() ContentEncoding
	SetContentEncoding(encoding ContentEncoding)
}

type aPart struct {
	*anObject
}

type rawPart interface {
	Object
	rawPart() *C.GMimePart
}

func CastPart(p *C.GMimePart) *aPart {
	return &aPart{CastObject((*C.GMimeObject)(unsafe.Pointer(p)))}
}

func NewPart() Part {
	part := C.g_mime_part_new()
	defer unref(C.gpointer(part))
	return CastPart(part)
}

func NewPartWithType(ctype string, csubtype string) Part {
	var _ctype *C.char = C.CString(ctype)
	var _csubtype *C.char = C.CString(csubtype)
	defer C.free(unsafe.Pointer(_ctype))
	defer C.free(unsafe.Pointer(_csubtype))

	part := C.g_mime_part_new_with_type(_ctype, _csubtype)
	defer unref(C.gpointer(part))
	return CastPart(part)
}

func (p *aPart) SetContentObject(content DataWrapper) {
	rawContent := content.(rawDataWrapper)
	C.g_mime_part_set_content_object(p.rawPart(), rawContent.rawDataWrapper())
}

func (p *aPart) ContentObject() DataWrapper {
	cDataWrapper := C.g_mime_part_get_content_object(p.rawPart())
	if cDataWrapper == nil {
		return nil
	}
	return CastDataWrapper(cDataWrapper)
}

func (p *aPart) Description() string {
	desc := C.g_mime_part_get_content_description(p.rawPart())
	return C.GoString(desc)
}

func (p *aPart) ContentLocation() string {
	loc := C.g_mime_part_get_content_location(p.rawPart())
	return C.GoString(loc)
}

func (p *aPart) ContentEncoding() ContentEncoding {
	return &aContentEncoding{encoding: C.g_mime_part_get_content_encoding(p.rawPart())}
}

func (p *aPart) SetContentEncoding(encoding ContentEncoding) {
	rawEncode := encoding.(rawContentEncoding)
	C.g_mime_part_set_content_encoding(p.rawPart(), rawEncode.rawContentEncoding())
}

func (p *aPart) Filename() string {
	return C.GoString(C.g_mime_part_get_filename(p.rawPart()))
}

func (p *aPart) rawPart() *C.GMimePart {
	return (*C.GMimePart)(p.pointer())
}
