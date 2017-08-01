package gmime

/*
#cgo pkg-config: gmime-3.0
#include <stdlib.h>
#include <strings.h>
#include <gmime/gmime.h>

static GMimeMessage *gmime_parse (const char *buffer, size_t len) {
	GMimeStream *stream = g_mime_stream_mem_new_with_buffer (buffer, len);
	GMimeParser *parser = g_mime_parser_new_with_stream (stream);
	g_object_unref (stream);
	GMimeMessage *message = g_mime_parser_construct_message (parser, NULL);
	g_object_unref (parser);

	InternetAddressList *list = g_mime_message_get_addresses (message, GMIME_ADDRESS_TYPE_TO);
	GMimeFormatOptions *format = g_mime_format_options_get_default ();
	char *buf = internet_address_list_to_string (list, format, FALSE);
	g_free (buf);

	int listLen = internet_address_list_length (list);
	for(int i = 0; i < listLen; i++) {
		InternetAddress *addr = internet_address_list_get_address (list, i);
		printf("Name: %s\n", internet_address_get_name (addr));
		printf("Address: %s\n", internet_address_mailbox_get_addr ((InternetAddressMailbox *)addr));
	}

	return message;
}

static char* gmime_get_content_type_string (GMimeObject *object) {
	GMimeContentType *ctype = g_mime_object_get_content_type (object);
	return g_mime_content_type_get_mime_type (ctype);
}

static char* gmime_get_content_string (GMimeObject *object) {
	if (!GMIME_IS_TEXT_PART (object)) {
		return NULL;
	}
	return g_mime_text_part_get_text ((GMimeTextPart *) object);
}

static gboolean gmime_is_text_part (GMimeObject *object) {
	return GMIME_IS_TEXT_PART (object);
}

static GByteArray *gmime_get_bytes (GMimeObject *object) {
	GMimeStream *stream;
	GMimeDataWrapper *content;
	GByteArray *buf;

	if (!(content = g_mime_part_get_content ((GMimePart *) object)))
		return NULL;
	stream = g_mime_stream_mem_new ();
	ssize_t size = g_mime_data_wrapper_write_to_stream (content, stream);
	printf("size: %zu\n", size);
	// g_mime_stream_flush (stream);

	buf = g_mime_stream_mem_get_byte_array ((GMimeStreamMem *) stream);
	g_mime_stream_mem_set_owner ((GMimeStreamMem *) stream, FALSE);

	g_object_unref (stream);
	return buf;
}

*/
import "C"
import "unsafe"
import "io"

// This function call automatically by runtime
func init() {
	C.g_mime_init()
}

// Shutdown is really needed only for valgrind
func Shutdown() {
	C.g_mime_shutdown()
}

// convert from Go bool to C gboolean
func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

// convert from C gboolean to Go bool
func gobool(b C.gboolean) bool {
	return b != C.gboolean(0)
}

// free up memory
func unref(referee C.gpointer) {
	C.g_object_unref(referee)
}

// Envelope wraps gmime message object and has methods to access it
type Envelope struct {
	gmimeMessage *C.GMimeMessage
}

// Part is a wrapper for message parts
type Part struct {
	gmimePart *C.GMimeObject
}

// ContentType returns part's content type
func (p *Part) ContentType() string {
	ctype := C.gmime_get_content_type_string(p.gmimePart)
	defer C.g_free(C.gpointer(unsafe.Pointer(ctype)))
	return C.GoString(ctype)
}

// IsText returns true if part's mime is text/*
func (p *Part) IsText() bool {
	return gobool(C.gmime_is_text_part(p.gmimePart))
}

// Text returns text portion of the part if it's mime is text/*
func (p *Part) Text() string {
	content := C.gmime_get_content_string(p.gmimePart)
	defer C.g_free(C.gpointer(unsafe.Pointer(content)))
	return C.GoString(content)
}

// Bytes returns decoded raw bytes of the part, most useful to access attachment data
func (p *Part) Bytes() []byte {
	b := C.gmime_get_bytes(p.gmimePart)
	defer C.g_byte_array_free((*C.GByteArray)(unsafe.Pointer(b)), C.TRUE)
	return C.GoBytes(unsafe.Pointer(b.data), C.int(b.len))
}

// SetText replaces text content if part is text/*
func (p *Part) SetText(text string) error {
	// TODO: Optimize this
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.g_mime_text_part_set_text((*C.GMimeTextPart)(unsafe.Pointer(p.gmimePart)), cstr)
	return nil
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
	stream := C.g_mime_stream_mem_new()                        // need unref
	defer C.g_object_unref(C.gpointer(unsafe.Pointer(stream))) // unref
	nWritten := C.g_mime_object_write_to_stream((*C.GMimeObject)(unsafe.Pointer(m.gmimeMessage)), nil, stream)
	if nWritten <= 0 {
		return nil, io.EOF
	}
	// byteArray is owned by stream and will be freed with it
	byteArray := C.g_mime_stream_mem_get_byte_array((*C.GMimeStreamMem)(unsafe.Pointer(stream)))
	return C.GoBytes(unsafe.Pointer(byteArray.data), (C.int)(nWritten)), nil
}

// Close frees up message resources
func (m *Envelope) Close() {
	C.g_object_unref(C.gpointer(m.gmimeMessage))
}
