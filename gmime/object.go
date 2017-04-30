package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
static gboolean object_is_part(GTypeInstance *obj) {
    return GMIME_IS_PART(obj);
}
static gboolean object_is_multipart(GTypeInstance *obj) {
    return GMIME_IS_MULTIPART(obj);
}
static gboolean object_is_message(GTypeInstance *obj) {
    return GMIME_IS_MESSAGE(obj);
}
static gboolean object_is_message_part(GTypeInstance *obj) {
    return GMIME_IS_MESSAGE_PART(obj);
}
static gboolean object_is_partial(GTypeInstance *obj) {
    return GMIME_IS_MESSAGE_PARTIAL(obj);
}
static GMimePart * gmime_part(GMimeObject *obj) {
	return GMIME_PART(obj);
}
*/
import "C"
import (
	"unsafe"
)

type Object interface {
	Janitor
	SetContentType(ContentType)
	ContentType() ContentType
	SetHeader(string, string)
	Header(string) (string, bool)
	ToString() string
	ContentDisposition() ContentDisposition
	Headers() string
	WriteToStream(Stream) int
	WalkHeaders(cb func(string, string) error) error
}

type anObject struct {
	*PointerMixin
}

type rawObject interface {
	Object
	rawObject() *C.GMimeObject
}

func CastObject(o *C.GMimeObject) *anObject {
	return &anObject{CastPointer(C.gpointer(o))}
}

func NewObject(contentType ContentType) Object {
	rawContentType := contentType.(rawContentType)
	object := C.g_mime_object_new(rawContentType.rawContentType())
	defer unref(C.gpointer(object))
	o := objectAsSubclass(object)
	return o
}

func NewObjectWithType(ctype string, csubtype string) Object {
	var _ctype *C.char = C.CString(ctype)
	var _csubtype *C.char = C.CString(csubtype)
	defer C.free(unsafe.Pointer(_ctype))
	defer C.free(unsafe.Pointer(_csubtype))

	object := C.g_mime_object_new_type(_ctype, _csubtype)
	defer unref(C.gpointer(object))
	o := objectAsSubclass(object)
	o.SetContentType(NewContentType(ctype, csubtype))

	return o
}

func (o *anObject) SetContentType(contentType ContentType) {
	rawContentType := contentType.(rawContentType)
	C.g_mime_object_set_content_type(o.rawObject(), rawContentType.rawContentType())
}

func (o *anObject) ContentType() ContentType {
	if ct := C.g_mime_object_get_content_type(o.rawObject()); ct != nil {
		return CastContentType(ct)
	}
	return nil
}

func (o *anObject) SetHeader(name, value string) {
	var _name *C.char = C.CString(name)
	var _value *C.char = C.CString(value)
	defer C.free(unsafe.Pointer(_name))
	defer C.free(unsafe.Pointer(_value))

	C.g_mime_object_set_header(o.rawObject(), _name, _value)
}

func (o *anObject) Header(name string) (string, bool) {
	var _name *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(_name))

	return maybeGoString(C.g_mime_object_get_header(o.rawObject(), _name))
}

func (o *anObject) ToString() string {
	str := C.g_mime_object_to_string(o.rawObject())
	defer C.free(unsafe.Pointer(str))
	return C.GoString(str)
}

func (o *anObject) ContentDisposition() ContentDisposition {
	cd := C.g_mime_object_get_content_disposition(o.rawObject())
	if cd != nil {
		return CastContentDisposition(cd)
	}
	return nil
}

func (o *anObject) Headers() string {
	headers := C.g_mime_object_get_headers(o.rawObject())
	defer C.free(unsafe.Pointer(headers))

	return C.GoString(headers)
}

func (o *anObject) WriteToStream(stream Stream) int {
	rawStream := stream.(rawStream)
	res := C.g_mime_object_write_to_stream(o.rawObject(), rawStream.rawStream())

	return int(res)
}

func (o *anObject) rawObject() *C.GMimeObject {
	return (*C.GMimeObject)(o.pointer())
}

func objectAsSubclass(o *C.GMimeObject) Object {
	partType := (*C.GTypeInstance)(unsafe.Pointer(o))

	if gobool(C.object_is_message_part(partType)) {
		return CastMessagePart((*C.GMimeMessagePart)(unsafe.Pointer(o)))
	} else if gobool(C.object_is_partial(partType)) {
		return CastMessagePartial((*C.GMimeMessagePartial)(unsafe.Pointer(o)))
	} else if gobool(C.object_is_multipart(partType)) {
		return CastMultipart((*C.GMimeMultipart)(unsafe.Pointer(o)))
	} else if gobool(C.object_is_part(partType)) {
		return CastPart((*C.GMimePart)(unsafe.Pointer(o)))
	} else if gobool(C.object_is_message(partType)) {
		return CastMessage((*C.GMimeMessage)(unsafe.Pointer(o)))
	} else {
		return CastObject(o)
	}
}

func (o *anObject) WalkHeaders(cb func(string, string) error) error {
	ghl := C.g_mime_object_get_header_list(o.rawObject())
	iter := C.g_mime_header_iter_new()
	defer C.g_mime_header_iter_free(iter)
	if !gobool(C.g_mime_header_list_get_iter(ghl, iter)) {
		return nil
	}
	for {
		name := C.GoString(C.g_mime_header_iter_get_name(iter))
		value := C.GoString(C.g_mime_header_iter_get_value(iter))
		err := cb(name, value)
		if err != nil {
			return err
		}
		if !gobool(C.g_mime_header_iter_next(iter)) {
			return nil
		}
	}
}

// Very minimal interface, to inspection only
type HeaderIterator interface {
	Janitor
	Name() string
	Value() string
	Next() bool
}

type aHeader struct {
}
