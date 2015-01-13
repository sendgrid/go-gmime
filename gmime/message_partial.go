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

type MessagePartial interface {
	Object
	Id() string
	Number() int
	Total() int
	//FIXME: temporary removed, wrong implemented
	//	ReconstructMessage(num int) Message
	//	SplitMessage(message Message, maxSize int, nParts *int) Message
}

type aMessagePartial struct {
	*aPart
}

func CastMessagePartial(mp *C.GMimeMessagePartial) MessagePartial {
	oPtr := (*C.GMimePart)(unsafe.Pointer(mp))
	return &aMessagePartial{CastPart(oPtr)}
}

func NewMessagePartial(id string, number int, total int) MessagePartial {
	cId := C.CString(id)
	defer C.free(unsafe.Pointer(cId))

	partial := C.g_mime_message_partial_new(cId, C.int(number), C.int(total))
	defer unref(C.gpointer(partial))
	return CastMessagePartial(partial)
}

func (m *aMessagePartial) Id() string {
	cId := C.g_mime_message_partial_get_id(m.rawPartial())

	return C.GoString(cId)
}

func (m *aMessagePartial) Number() int {
	cNumber := C.g_mime_message_partial_get_number(m.rawPartial())

	return int(cNumber)
}

func (m *aMessagePartial) Total() int {
	cTotal := C.g_mime_message_partial_get_total(m.rawPartial())

	return int(cTotal)
}

/* WRONG (and possibly unneeded)
func (m *aMessagePartial) ReconstructMessage(num int) Message {
	cMessage := C.g_mime_message_partial_reconstruct_message(m.rawPartial(), C.size_t(num))

    defer unref(C.gpointer(cMessage))
    return CastMessage(cMessage)
}

*/

/* WRONG (and possibly unneeded)
func (m *aMessagePartial) SplitMessage(message Message, maxSize int, nParts *int) Message {
	rawMessage := message.(rawMessage)
	cNparts := (*C.size_t)(unsafe.Pointer(nParts))
	cMessages := C.g_mime_message_partial_split_message(rawMessage.rawMessage(), C.size_t(maxSize), cNparts)

	return &aMessage{
		messages: cMessages,
	}
}
*/

func (mp *aMessagePartial) rawPartial() *C.GMimeMessagePartial {
	return (*C.GMimeMessagePartial)(mp.pointer())
}
