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

type Message interface {
	Object
	SetSender(string)
	Sender() (string, bool)
	SetReplyTo(string)
	ReplyTo() (string, bool)
	SetSubject(string)
	Subject() (string, bool)
	SetMessageId(string)
	MessageId() (string, bool)
	DateAsString() (string, bool)
	SetDateAsString(string)
	AddTo(string, string)
	To() *InternetAddressList
	AddCc(string, string)
	Cc() *InternetAddressList
	AddBcc(string, string)
	Bcc() *InternetAddressList
	AllRecipients() *InternetAddressList
	SetMimePart(Object)
	MimePart() Object
	Body() Object
	Date() string
}

type aMessage struct {
	*anObject
}

type rawMessage interface {
	Message
	rawMessage() *C.GMimeMessage
}

func CastMessage(m *C.GMimeMessage) *aMessage {
	return &aMessage{CastObject((*C.GMimeObject)(unsafe.Pointer(m)))}
}

func NewMessage() Message {
	m := C.g_mime_message_new(0)
	defer unref(C.gpointer(m))
	return CastMessage(m)
}

func (m *aMessage) SetSender(sender string) {
	var cSender *C.char = C.CString(sender)
	C.g_mime_message_set_sender(m.rawMessage(), cSender)
	C.free(unsafe.Pointer(cSender))
}

func (m *aMessage) Sender() (string, bool) {
	sender := C.g_mime_message_get_sender(m.rawMessage())
	return maybeGoString(sender)
}

func (m *aMessage) SetReplyTo(replyTo string) {
	var cReply *C.char = C.CString(replyTo)
	C.g_mime_message_set_reply_to(m.rawMessage(), cReply)
	C.free(unsafe.Pointer(cReply))
}

func (m *aMessage) ReplyTo() (string, bool) {
	replyTo := C.g_mime_message_get_reply_to(m.rawMessage())
	return maybeGoString(replyTo)
}

func (m *aMessage) SetSubject(subject string) {
	var cSubject *C.char = C.CString(subject)
	C.g_mime_message_set_subject(m.rawMessage(), cSubject)
	C.free(unsafe.Pointer(cSubject))
}

func (m *aMessage) Subject() (string, bool) {
	subject := C.g_mime_message_get_subject(m.rawMessage())
	return maybeGoString(subject)
}

func (m *aMessage) SetMessageId(messageId string) {
	var cMessageId *C.char = C.CString(messageId)
	C.g_mime_message_set_message_id(m.rawMessage(), cMessageId)
	C.free(unsafe.Pointer(cMessageId))
}

func (m *aMessage) MessageId() (string, bool) {
	messageId := C.g_mime_message_get_message_id(m.rawMessage())
	return maybeGoString(messageId)
}

func (m *aMessage) DateAsString() (string, bool) {
	date := C.g_mime_message_get_date_as_string(m.rawMessage())
	defer C.free(unsafe.Pointer(date))
	return maybeGoString(date)
}

func (m *aMessage) SetDateAsString(date string) {
	var cDate *C.char = C.CString(date)
	C.g_mime_message_set_date_as_string(m.rawMessage(), cDate)
	C.free(unsafe.Pointer(cDate))
}

func (m *aMessage) addRecipient(recipientType C.GMimeRecipientType, name string, address string) {
	var cName *C.char = C.CString(name)
	var cAddress *C.char = C.CString(address)
	C.g_mime_message_add_recipient(m.rawMessage(), recipientType, cName, cAddress)
	C.free(unsafe.Pointer(cName))
	C.free(unsafe.Pointer(cAddress))
}

func (m *aMessage) AddTo(name string, address string) {
	m.addRecipient(C.GMIME_RECIPIENT_TYPE_TO, name, address)
}

func (m *aMessage) To() *InternetAddressList {
	cList := C.g_mime_message_get_recipients(m.rawMessage(), C.GMIME_RECIPIENT_TYPE_TO)
	return CastInternetAddressList(cList)
}

func (m *aMessage) AddCc(name string, address string) {
	m.addRecipient(C.GMIME_RECIPIENT_TYPE_CC, name, address)
}

func (m *aMessage) Cc() *InternetAddressList {
	cList := C.g_mime_message_get_recipients(m.rawMessage(), C.GMIME_RECIPIENT_TYPE_CC)
	return CastInternetAddressList(cList)
}

func (m *aMessage) AddBcc(name string, address string) {
	m.addRecipient(C.GMIME_RECIPIENT_TYPE_BCC, name, address)
}

func (m *aMessage) Bcc() *InternetAddressList {
	cList := C.g_mime_message_get_recipients(m.rawMessage(), C.GMIME_RECIPIENT_TYPE_BCC)
	return CastInternetAddressList(cList)
}

func (m *aMessage) AllRecipients() *InternetAddressList {
	// This is major exception: we have newly allocated list here
	cList := C.g_mime_message_get_all_recipients(m.rawMessage())
	defer unref(C.gpointer(cList))
	return CastInternetAddressList(cList)
}

func (m *aMessage) SetMimePart(mimePart Object) {
	part := mimePart.(rawObject)
	switch mimePart.(type) {
	case Part:
		C.g_mime_message_set_mime_part(m.rawMessage(), part.rawObject())
	case Multipart:
		C.g_mime_message_set_mime_part(m.rawMessage(), part.rawObject())
	}
}

func (m *aMessage) MimePart() Object {
	object := C.g_mime_message_get_mime_part(m.rawMessage())
	return objectAsSubclass(object)
}

func (m *aMessage) Body() Object {
	object := C.g_mime_message_get_body(m.rawMessage())
	return objectAsSubclass(object)
}

func (m *aMessage) Date() string {
	cDate := C.g_mime_message_get_date_as_string(m.rawMessage())

	return C.GoString(cDate)
}

func (m *aMessage) rawMessage() *C.GMimeMessage {
	return (*C.GMimeMessage)(m.pointer())
}
