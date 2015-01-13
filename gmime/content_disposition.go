package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
static gboolean object_is_content_disposition(GMimeContentDisposition *obj) {
	return GMIME_IS_CONTENT_DISPOSITION(obj);
}
*/
import "C"
import (
	"unsafe"
)

type ContentDisposition interface {
	Janitor
	Parametrized
	Disposition() string
	ToString(fold bool) string
	IsAttachment() bool
}

type aContentDisposition struct {
	*PointerMixin
}

func CastContentDisposition(c *C.GMimeContentDisposition) *aContentDisposition {
	return &aContentDisposition{CastPointer(C.gpointer(c))}
}

func NewContentDispositionFromString(str string) *aContentDisposition {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	cd := C.g_mime_content_disposition_new_from_string(cStr)
	defer unref(C.gpointer(cd))

	return CastContentDisposition(cd)
}

func (d *aContentDisposition) Disposition() string {
	if gobool(C.object_is_content_disposition(d.rawContentDisposition())) {
		cDisposition := C.g_mime_content_disposition_get_disposition(d.rawContentDisposition())

		return C.GoString(cDisposition)
	}
	return ""
}

func (d *aContentDisposition) IsAttachment() bool {
	if d.Disposition() != "" {
		// it should be either "attachment" or "inline"
		return true
	}

	return false
}

func (d *aContentDisposition) ToString(fold bool) string {
	cDisposition := C.g_mime_content_disposition_to_string(d.rawContentDisposition(), gbool(fold))

	defer C.free(unsafe.Pointer(cDisposition))
	return C.GoString(cDisposition)
}

func (t *aContentDisposition) SetParameter(name, value string) {
	var _name *C.char = C.CString(name)
	var _value *C.char = C.CString(value)
	defer C.free(unsafe.Pointer(_name))
	defer C.free(unsafe.Pointer(_value))
	C.g_mime_content_disposition_set_parameter(t.rawContentDisposition(), _name, _value)
}

func (t *aContentDisposition) Parameter(name string) string {
	var _name *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(_name))
	content := C.g_mime_content_disposition_get_parameter(t.rawContentDisposition(), _name)
	return C.GoString(content)
}

func (t *aContentDisposition) ForEachParam(callback GMimeParamsCallback) {
	params := C.g_mime_content_disposition_get_params(t.rawContentDisposition())
	forEachParam(params, callback)
}

func (d *aContentDisposition) rawContentDisposition() *C.GMimeContentDisposition {
	return ((*C.GMimeContentDisposition)(d.pointer()))
}
