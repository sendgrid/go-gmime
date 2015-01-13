package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"
import "unsafe"

type MessagePart interface {
	Object
	Message() Message
	SetMessage(message Message)
}

type aMessagePart struct {
	*anObject
}

type rawMessagePart interface {
	MessagePart
	rawMessagePart() *C.GMimeMessagePart
}

func CastMessagePart(mp *C.GMimeMessagePart) MessagePart {
	oPtr := (*C.GMimeObject)(unsafe.Pointer(mp))
	return &aMessagePart{CastObject(oPtr)}
}

func NewMessagePart(subtype string) MessagePart {
	var cSubtype *C.char = C.CString(subtype)
	defer C.free(unsafe.Pointer(cSubtype))

	part := C.g_mime_message_part_new(cSubtype)
	defer unref(C.gpointer(part))
	return CastMessagePart(part)
}

func NewMessagePartWithMessage(subtype string, message Message) MessagePart {
	rawMessage := message.(rawMessage)
	var cSubtype *C.char = C.CString(subtype)
	defer C.free(unsafe.Pointer(cSubtype))

	part := C.g_mime_message_part_new_with_message(cSubtype, rawMessage.rawMessage())
	defer unref(C.gpointer(part))
	return CastMessagePart(part)
}

func (m *aMessagePart) Message() Message {
	message := C.g_mime_message_part_get_message(m.rawMessagePart())
	return CastMessage(message)
}

func (m *aMessagePart) SetMessage(message Message) {
	rawMessage := message.(rawMessage)
	C.g_mime_message_part_set_message(m.rawMessagePart(), rawMessage.rawMessage())
}

func (m *aMessagePart) rawMessagePart() *C.GMimeMessagePart {
	return (*C.GMimeMessagePart)(m.pointer())
}
