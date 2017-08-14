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
	C.g_mime_message_set_subject(m.gmimeMessage, C.CString(subject), C.CString("UTF-8"))
}

// Headers returns all headers for envelope
func (m *Envelope) Headers() textproto.MIMEHeader {
	// TODO: this is not super efficient, but easier to read, may be optimize this?
	headers := C.g_mime_object_get_header_list(m.asGMimeObject())
	count := C.g_mime_header_list_get_count(headers)
	goHeaders := make(textproto.MIMEHeader, int(count))
	var i C.int
	for i = 0; i < count; i++ {
		header := C.g_mime_header_list_get_header_at(headers, i)
		name := C.GoString(C.g_mime_header_get_name(header))
		value := C.GoString(C.g_mime_header_get_value(header))
		if _, ok := goHeaders[name]; !ok {
			goHeaders[name] = nil
		}
		goHeaders[name] = append(goHeaders[name], value)
	}
	return goHeaders
}

// SetHeader sets or replaces specified header
func (m *Envelope) SetHeader(name string, value string) {
	headers := C.g_mime_object_get_header_list(m.asGMimeObject())
	C.g_mime_header_list_set(headers, C.CString(name), C.CString(value), C.CString("UTF-8"))
}

// RemoveHeader removes existing header
func (m *Envelope) RemoveHeader(name string) bool {
	headers := C.g_mime_object_get_header_list(m.asGMimeObject())
	return gobool(C.g_mime_header_list_remove(headers, C.CString(name)))
}

// Header returns *first* header from envelope
// if user wants to get all headers use `Headers` function
func (m *Envelope) Header(header string) string {
	cHeaderName := C.CString(header)
	defer C.free(unsafe.Pointer(cHeaderName))
	cHeader := C.g_mime_object_get_header(m.asGMimeObject(), cHeaderName)
	return C.GoString(cHeader)
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
	partIter := C.g_mime_part_iter_new(m.asGMimeObject())
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
	format := C.g_mime_format_options_get_default()
	C.g_mime_format_options_set_newline_format(format, C.GMIME_NEWLINE_FORMAT_DOS)
	stream := C.g_mime_stream_mem_new()                        // need unref
	defer C.g_object_unref(C.gpointer(unsafe.Pointer(stream))) // unref
	nWritten := C.g_mime_object_write_to_stream(m.asGMimeObject(), format, stream)
	if nWritten <= 0 {
		return nil, errors.New("can't write to stream")
	}
	// byteArray is owned by stream and will be freed with it
	byteArray := C.g_mime_stream_mem_get_byte_array((*C.GMimeStreamMem)(unsafe.Pointer(stream)))
	return C.GoBytes(unsafe.Pointer(byteArray.data), (C.int)(byteArray.len)), nil
}

// Close frees up message resources
func (m *Envelope) Close() {
	C.g_object_unref(C.gpointer(m.gmimeMessage))
}

func (m *Envelope) asGMimeObject() *C.GMimeObject {
	return (*C.GMimeObject)(unsafe.Pointer(m.gmimeMessage))
}
