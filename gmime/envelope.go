package gmime

// #include "gmime.h"
import "C"
import (
	"errors"
	"fmt"
	"net/mail"
	"net/textproto"
	"strings"
	"unsafe"
)

// Envelope wraps gmime message object and has methods to access it
type Envelope struct {
	gmimeMessage *C.GMimeMessage
}

// Parse parses message and returns Message
func Parse(data string) (*Envelope, error) {
	// very inefficient
	cBuf := C.CString(data)
	defer C.free(unsafe.Pointer(cBuf))
	gmsg := C.gmime_parse(cBuf, C.size_t(len(data)))
	if gmsg == nil {
		return nil, fmt.Errorf("gmime.parse: unable to parse mime")
	}

	return &Envelope{
		gmimeMessage: gmsg,
	}, nil
}

// Subject returns envelope's Subject
// gmime returns the string in utf-8
func (m *Envelope) Subject() string {
	subject := C.g_mime_message_get_subject(m.gmimeMessage)
	return C.GoString(subject)
}

// SetSubject returns envelope's Subject
func (m *Envelope) SetSubject(subject string) {
	cSubject := C.CString(subject)
	defer C.free(unsafe.Pointer(cSubject))
	cType := C.CString("UTF-8")
	defer C.free(unsafe.Pointer(cType))
	C.g_mime_message_set_subject(m.gmimeMessage, cSubject, cType)
}

// Headers returns all headers for envelope
func (m *Envelope) Headers() textproto.MIMEHeader {
	// TODO: this is not super efficient, but easier to read, may be optimize this?
	headers := C.g_mime_object_get_header_list(m.asGMimeObject())
	count := C.g_mime_header_list_get_count(headers)
	goHeaders := make(textproto.MIMEHeader)
	var i C.int
	for i = 0; i < count; i++ {
		header := C.g_mime_header_list_get_header_at(headers, i)
		name := C.GoString(C.g_mime_header_get_name(header))
		value := C.GoString(C.g_mime_header_get_value(header))
		goHeaders.Add(name, value)
	}
	return goHeaders
}

// RawHeaders returns all headers for envelope without rfc2047 decoding
func (m *Envelope) RawHeaders() textproto.MIMEHeader {
	// TODO: this is not super efficient, but easier to read, may be optimize this?
	headers := C.g_mime_object_get_header_list(m.asGMimeObject())
	count := C.g_mime_header_list_get_count(headers)
	goHeaders := make(textproto.MIMEHeader)
	var i C.int
	for i = 0; i < count; i++ {
		header := C.g_mime_header_list_get_header_at(headers, i)
		name := C.GoString(C.g_mime_header_get_name(header))
		value := C.GoString(C.g_mime_header_get_raw_value(header))
		goHeaders.Add(name, value)
	}
	return goHeaders
}

// ReplaceHeader replaces the header with matching key & originalValue with replaceValue
func (m *Envelope) ReplaceHeader(key, originalValue, replaceValue string) error {
	headers := C.g_mime_object_get_header_list(m.asGMimeObject())
	count := C.g_mime_header_list_get_count(headers)
	var i C.int
	for i = 0; i < count; i++ {
		header := C.g_mime_header_list_get_header_at(headers, i)
		name := C.GoString(C.g_mime_header_get_name(header))
		if name != key {
			continue
		}
		value := C.GoString(C.g_mime_header_get_value(header))
		if value != originalValue {
			continue
		}
		format := C.g_mime_format_options_get_default()
		cCharset := C.CString("UTF-8")
		cReplaceValue := C.CString(replaceValue)
		C.g_mime_header_set_value(header, format, cReplaceValue, cCharset)
		C.free(unsafe.Pointer(cReplaceValue))
		C.free(unsafe.Pointer(cCharset))
		return nil
	}
	return errors.New(fmt.Sprintf("failed to find header with matching key: %v & value: %v", key, originalValue))
}

// SetHeader sets or replaces specified header
func (m *Envelope) SetHeader(name string, value string) error {
	switch strings.ToLower(name) {
	case "from", "sender", "reply-to", "to", "cc", "bcc":
		return fmt.Errorf("use AddAddress for %s", name)
	default:
		headers := C.g_mime_object_get_header_list(m.asGMimeObject())
		cName := C.CString(name)
		defer C.free(unsafe.Pointer(cName))
		cValue := C.CString(value)
		defer C.free(unsafe.Pointer(cValue))
		cCharset := C.CString("UTF-8")
		defer C.free(unsafe.Pointer(cCharset))

		C.g_mime_header_list_set(headers, cName, cValue, cCharset)
		return nil
	}
}

// AddAddress adds an address from/sender/reply-to/to to/cc/bcc
func (m *Envelope) AddAddress(header, name, address string) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cAddress := C.CString(address)
	defer C.free(unsafe.Pointer(cAddress))

	var addressList *C.InternetAddressList
	switch strings.ToLower(header) {
	case "from":
		addressList = C.g_mime_message_get_from(m.gmimeMessage)
	case "sender":
		addressList = C.g_mime_message_get_sender(m.gmimeMessage)
	case "reply-to":
		addressList = C.g_mime_message_get_reply_to(m.gmimeMessage)
	case "to":
		addressList = C.g_mime_message_get_to(m.gmimeMessage)
	case "cc":
		addressList = C.g_mime_message_get_cc(m.gmimeMessage)
	case "bcc":
		addressList = C.g_mime_message_get_bcc(m.gmimeMessage)
	default:
		return fmt.Errorf("can't add to header %s", header)
	}

	mb := C.internet_address_mailbox_new(cName, cAddress)
	C.internet_address_list_add(addressList, mb)
	return nil
}

// ParseAndAppendAddresses will attempt to parse a string and appends to the list
func (m *Envelope) ParseAndAppendAddresses(header, addresses string) error {
	cAddresses := C.CString(addresses)
	defer C.free(unsafe.Pointer(cAddresses))

	var addressList *C.InternetAddressList
	switch strings.ToLower(header) {
	case "from":
		addressList = C.g_mime_message_get_from(m.gmimeMessage)
	case "sender":
		addressList = C.g_mime_message_get_sender(m.gmimeMessage)
	case "reply-to":
		addressList = C.g_mime_message_get_reply_to(m.gmimeMessage)
	case "to":
		addressList = C.g_mime_message_get_to(m.gmimeMessage)
	case "cc":
		addressList = C.g_mime_message_get_cc(m.gmimeMessage)
	case "bcc":
		addressList = C.g_mime_message_get_bcc(m.gmimeMessage)
	default:
		return fmt.Errorf("can't append addresses to header %s", header)
	}

	parsed := C.internet_address_list_parse(C.g_mime_parser_options_get_default(), cAddresses)
	if parsed != nil {
		defer C.g_object_unref((C.gpointer)(unsafe.Pointer(parsed)))
		C.internet_address_list_append(addressList, parsed)
	}
	return nil
}

// AppendAddressList appends a list of mail.Addresses to the specified header
func (m *Envelope) AppendAddressList(header string, addrs []*mail.Address) error {
	var addressList *C.InternetAddressList
	switch strings.ToLower(header) {
	case "from":
		addressList = C.g_mime_message_get_from(m.gmimeMessage)
	case "sender":
		addressList = C.g_mime_message_get_sender(m.gmimeMessage)
	case "reply-to":
		addressList = C.g_mime_message_get_reply_to(m.gmimeMessage)
	case "to":
		addressList = C.g_mime_message_get_to(m.gmimeMessage)
	case "cc":
		addressList = C.g_mime_message_get_cc(m.gmimeMessage)
	case "bcc":
		addressList = C.g_mime_message_get_bcc(m.gmimeMessage)
	default:
		return fmt.Errorf("can't append addresses to header %s", header)
	}
	for _, addr := range addrs {
		name := C.CString(addr.Name)
		addr := C.CString(addr.Address)
		mb := C.internet_address_mailbox_new(name, addr)
		C.free(unsafe.Pointer(name))
		C.free(unsafe.Pointer(addr))
		C.internet_address_list_add(addressList, mb)
	}
	return nil
}

// ClearAddress will clear the from/sender/reply-to/to/cc/bcc list
// however, it will not clear the header name!
// if you want the entire header removed, use RemoveHeader
func (m *Envelope) ClearAddress(header string) error {
	var addressList *C.InternetAddressList
	switch strings.ToLower(header) {
	case "from":
		addressList = C.g_mime_message_get_from(m.gmimeMessage)
	case "sender":
		addressList = C.g_mime_message_get_sender(m.gmimeMessage)
	case "reply-to":
		addressList = C.g_mime_message_get_reply_to(m.gmimeMessage)
	case "to":
		addressList = C.g_mime_message_get_to(m.gmimeMessage)
	case "cc":
		addressList = C.g_mime_message_get_cc(m.gmimeMessage)
	case "bcc":
		addressList = C.g_mime_message_get_bcc(m.gmimeMessage)
	default:
		return fmt.Errorf("unknown header %s", header)
	}

	C.internet_address_list_clear(addressList)
	return nil
}

// RemoveHeader removes existing header
func (m *Envelope) RemoveHeader(name string) bool {
	headers := C.g_mime_object_get_header_list(m.asGMimeObject())
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return gobool(C.g_mime_header_list_remove(headers, cName))
}

// RemoveAllHeaders removes all headers with the name if there are multiple
func (m *Envelope) RemoveAllHeaders(name string) bool {
	headers := C.g_mime_object_get_header_list(m.asGMimeObject())
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	next := true
	removed := false
	for next {
		next = C.g_mime_header_list_remove(headers, cName) == 1
		removed = removed || next
	}

	return removed
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
	if mimePart != nil {
		ctype := C.gmime_get_content_type_string(mimePart)
		defer C.g_free(C.gpointer(unsafe.Pointer(ctype)))
		return C.GoString(ctype)
	}
	return ""
}

func (m *Envelope) ContentTypeWithParameters() string {
	mimePart := C.g_mime_message_get_mime_part(m.gmimeMessage)
	if mimePart != nil {
		ctype := C.gmime_get_content_string_full(mimePart)
		defer C.g_free(C.gpointer(unsafe.Pointer(ctype)))
		return strings.TrimSpace(C.GoString(ctype))
	}
	return ""
}

// Walk iterates all message parts and executes callback on each part
func (m *Envelope) Walk(cb func(p *Part) error) error {
	partIter := C.g_mime_part_iter_new(m.asGMimeObject())
	defer C.g_mime_part_iter_free(partIter)
	for {
		currentPart := C.g_mime_part_iter_get_current(partIter)
		part := &Part{
			gmimePart: currentPart,
		}
		err := cb(part)
		if err != nil {
			return err
		}
		next := C.g_mime_part_iter_next(partIter)
		if !gobool(next) {
			break
		}
	}
	return nil
}

func (m *Envelope) WalkWithParent(cb func(p *Part) error) error {
	partIter := C.g_mime_part_iter_new(m.asGMimeObject())
	defer C.g_mime_part_iter_free(partIter)

	for {
		currentPart := C.g_mime_part_iter_get_current(partIter)

		parentPartObj := C.g_mime_part_iter_get_parent(partIter)
		var parentPart *Part
		if parentPartObj != nil {
			parentPart = &Part{
				gmimePart: parentPartObj,
			}
		}

		part := &Part{
			gmimePart: currentPart,
			parent:    parentPart,
		}

		err := cb(part)
		if err != nil {
			return err
		}

		next := C.g_mime_part_iter_next(partIter)
		if !gobool(next) {
			break
		}
	}

	return nil
}

// Export composes mime from envelope
func (m *Envelope) Export() ([]byte, error) {
	// TODO: optimize this, bundle cgo calls
	stream := C.g_mime_stream_mem_new()                        // need unref
	defer C.g_object_unref(C.gpointer(unsafe.Pointer(stream))) // unref
	format := C.g_mime_format_options_get_default()
	C.g_mime_format_options_set_newline_format(format, C.GMIME_NEWLINE_FORMAT_DOS)
	nWritten := C.g_mime_object_write_to_stream(m.asGMimeObject(), format, stream)
	if nWritten <= 0 {
		return nil, errors.New("can't write to stream")
	}
	// byteArray is owned by stream and will be freed with it
	// TODO: we can optimize it, avoiding copy, but will have to free manually
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

// AddHTMLAlternativeToPlainText converts plain text to html
func (m *Envelope) AddHTMLAlternativeToPlainText(content string) bool {
	rootPart := C.g_mime_message_get_mime_part(m.gmimeMessage)
	ctype := C.gmime_get_content_type_string(rootPart)
	defer C.g_free(C.gpointer(unsafe.Pointer(ctype)))

	// if content type is not text/plain, we can't reliably add alternative html
	goCType := C.GoString(ctype)
	if goCType != "text/plain" {
		return false
	}

	// set new multipart as root and add the original text/plain to it
	multipart := C.g_mime_multipart_new_with_subtype(cStringAlternative)
	C.g_mime_multipart_add(multipart, (*C.GMimeObject)(unsafe.Pointer(rootPart)))
	C.g_mime_message_set_mime_part(m.gmimeMessage, (*C.GMimeObject)(unsafe.Pointer(multipart)))

	// create a new html part and add it to the root level multipart
	// TODO: need to free newHTMLpart somewhere
	newHTMLpart := C.g_mime_text_part_new_with_subtype(cStringHTML)
	cContent := C.CString(content)
	defer C.g_free(C.gpointer(unsafe.Pointer(cContent)))
	C.g_mime_text_part_set_text(newHTMLpart, cContent)
	C.g_mime_multipart_add(multipart, (*C.GMimeObject)(unsafe.Pointer(newHTMLpart)))
	return true
}
