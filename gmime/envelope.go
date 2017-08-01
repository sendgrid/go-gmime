package gmime

/*
#cgo pkg-config: gmime-3.0
#include "gmime.h"
*/
import "C"

import (
	"errors"
	"net/textproto"
	"unsafe"
)

// Envelope wraps gmime message object and has methods to access it
type Envelope struct {
	gmimeMessage *C.GMimeMessage
}

// Parse parses message and returns Message
func Parse(data string) *Envelope {
	// very inefficient
	cBuf := C.CString(data)
	defer C.free(unsafe.Pointer(cBuf))
	return &Envelope{
		gmimeMessage: C.gmime_parse(cBuf, C.size_t(len(data))),
	}
}

// Subject returns envelope's Subject
func (m *Envelope) Subject() string {
	subject := C.g_mime_message_get_subject(m.gmimeMessage)
	return C.GoString(subject)
}

// SetSubject returns envelope's Subject
func (m *Envelope) SetSubject(subject string) {
}

func (m *Envelope) Headers() textproto.MIMEHeader {
	C.gmime_get_headers((*C.GMimeObject)(unsafe.Pointer(m.gmimeMessage)))
	return nil
}

func (m *Envelope) SetHeader() textproto.MIMEHeader {
	return nil
}

func (m *Envelope) Header(header string) []string {
	cHeaderName := C.CString(header)
	defer C.free(unsafe.Pointer(cHeaderName))
	cHeader := C.g_mime_object_get_header((*C.GMimeObject)(unsafe.Pointer(m.gmimeMessage)), cHeaderName)
	return []string{C.GoString(cHeader)}
}

// ContentType returns envelope's content-type
func (m *Envelope) ContentType() string {
	mimePart := C.g_mime_message_get_mime_part(m.gmimeMessage)
	ctype := C.gmime_get_content_type_string(mimePart)
	defer C.g_free(C.gpointer(unsafe.Pointer(ctype)))
	return C.GoString(ctype)
}

// Walk iterates all message parts and executes callback on each part
func (m *Envelope) Walk(cb func(p *Part)) {
	partIter := C.g_mime_part_iter_new((*C.GMimeObject)(unsafe.Pointer(m.gmimeMessage)))
	for {
		next := C.g_mime_part_iter_next(partIter)
		if !gobool(next) {
			break
		}
		currentPart := C.g_mime_part_iter_get_current(partIter)
		part := &Part{
			gmimePart: currentPart,
		}
		cb(part)
	}
}

// Export composes mime from envelope
func (m *Envelope) Export() ([]byte, error) {
	// TODO: optimize this, bundle cgo calls
	stream := C.g_mime_stream_mem_new()                        // need unref
	defer C.g_object_unref(C.gpointer(unsafe.Pointer(stream))) // unref
	nWritten := C.g_mime_object_write_to_stream((*C.GMimeObject)(unsafe.Pointer(m.gmimeMessage)), nil, stream)
	if nWritten <= 0 {
		return nil, errors.New("can't write to stream")
	}
	// byteArray is owned by stream and will be freed with it
	byteArray := C.g_mime_stream_mem_get_byte_array((*C.GMimeStreamMem)(unsafe.Pointer(stream)))
	return C.GoBytes(unsafe.Pointer(byteArray.data), (C.int)(nWritten)), nil
}

// Close frees up message resources
func (m *Envelope) Close() {
	C.g_object_unref(C.gpointer(m.gmimeMessage))
}
