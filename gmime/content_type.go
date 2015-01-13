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

type ContentType interface {
	Janitor
	Parametrized
	ToString() string
	MediaType() string
	MediaSubtype() string
}

type rawContentType interface {
	ContentType
	rawContentType() *C.GMimeContentType
}

type aContentType struct {
	*PointerMixin
}

func CastContentType(ct *C.GMimeContentType) ContentType {
	return &aContentType{CastPointer(C.gpointer(ct))}
}

func NewContentType(ctype string, csubtype string) ContentType {
	var _ctype *C.char = C.CString(ctype)
	var _csubtype *C.char = C.CString(csubtype)
	defer C.free(unsafe.Pointer(_ctype))
	defer C.free(unsafe.Pointer(_csubtype))

	ct := C.g_mime_content_type_new(_ctype, _csubtype)
	defer unref(C.gpointer(ct))

	return CastContentType(ct)
}

func NewContentTypeFromString(str string) ContentType {
	var _str *C.char = C.CString(str)
	defer C.free(unsafe.Pointer(_str))

	ct := C.g_mime_content_type_new_from_string(_str)
	defer unref(C.gpointer(ct))

	return CastContentType(ct)
}

func (t *aContentType) ToString() string {
	var contentType *C.char = C.g_mime_content_type_to_string(t.rawContentType())
	defer C.free(unsafe.Pointer(contentType))

	return C.GoString(contentType)
}

func (t *aContentType) MediaType() string {
	var mediaType *C.char = C.g_mime_content_type_get_media_type(t.rawContentType())

	return C.GoString(mediaType)
}

func (t *aContentType) MediaSubtype() string {
	var mediaSubtype *C.char = C.g_mime_content_type_get_media_subtype(t.rawContentType())

	return C.GoString(mediaSubtype)
}

func (t *aContentType) SetParameter(name, value string) {
	var _name *C.char = C.CString(name)
	var _value *C.char = C.CString(value)
	defer C.free(unsafe.Pointer(_name))
	defer C.free(unsafe.Pointer(_value))
	C.g_mime_content_type_set_parameter(t.rawContentType(), _name, _value)
}

func (t *aContentType) Parameter(name string) string {
	var _name *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(_name))
	content := C.g_mime_content_type_get_parameter(t.rawContentType(), _name)
	return C.GoString(content)
}

func (t *aContentType) ForEachParam(callback GMimeParamsCallback) {
	params := C.g_mime_content_type_get_params(t.rawContentType())
	forEachParam(params, callback)
}

func (t *aContentType) rawContentType() *C.GMimeContentType {
	return ((*C.GMimeContentType)(t.pointer()))
}
