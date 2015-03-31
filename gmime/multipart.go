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

type Multipart interface {
	Object
	AddPart(Object)
	Boundary() string
	GetPart(int) Object
	Count() int
	Clear()
    Walk(func (Object) error) error
}

type aMultipart struct {
	*anObject
}

type rawMultipart interface {
	Multipart
	rawMultipart() *C.GMimeMultipart
}

func CastMultipart(mp *C.GMimeMultipart) Multipart {
	oPtr := (*C.GMimeObject)(unsafe.Pointer(mp))
	return &aMultipart{CastObject(oPtr)}
}

func NewMultipart() Multipart {
	multipart := C.g_mime_multipart_new()
	defer unref(C.gpointer(multipart))
	return CastMultipart(multipart)
}

func NewMultipartWithSubtype(subtype string) Multipart {
	var _csubtype *C.char = C.CString(subtype)
	defer C.free(unsafe.Pointer(_csubtype))
	multipart := C.g_mime_multipart_new_with_subtype(_csubtype)
	defer unref(C.gpointer(multipart))
	return CastMultipart(multipart)
}

func (m *aMultipart) AddPart(part Object) {
	rawObject := part.(rawObject)
	C.g_mime_multipart_add(m.rawMultipart(), rawObject.rawObject())
}

func (m *aMultipart) GetPart(index int) Object {
	part := C.g_mime_multipart_get_part(m.rawMultipart(), C.int(index))
	return objectAsSubclass(part)
}

func (m *aMultipart) Clear() {
	C.g_mime_multipart_clear(m.rawMultipart())
}

func (m *aMultipart) Count() int {
	return int(C.g_mime_multipart_get_count(m.rawMultipart()))
}

func (m *aMultipart) Boundary() string {
	return C.GoString(C.g_mime_multipart_get_boundary(m.rawMultipart()))
}

// FIXME: need tests
func (m *aMultipart) Walk(callback func(Object) error) error {
	for i := 0; i < m.Count(); i++ {
		if ok := callback(m.GetPart(i)); ok != nil {
			return ok
		}
	}
	return nil
}

func (m *aMultipart) rawMultipart() *C.GMimeMultipart {
	return (*C.GMimeMultipart)(m.pointer())
}
