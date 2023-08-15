package gmime

// #include "gmime.h"
import "C"
import (
	"regexp"
	"unsafe"
)

import (
	"net/textproto"
	"strings"
)

// Part is a wrapper for message parts
type Part struct {
	gmimePart *C.GMimeObject
	parent    *Part
}

// ContentType returns part's content type
func (p *Part) ContentType() string {
	ctype := C.gmime_get_content_type_string(p.gmimePart)
	defer C.g_free(C.gpointer(unsafe.Pointer(ctype)))
	return C.GoString(ctype)
}

// ContentTypeWithParam returns content type's parameter
func (p *Part) ContentTypeWithParam(param string) string {
	ctype := C.g_mime_object_get_content_type(p.gmimePart)
	charset := C.g_mime_content_type_get_parameter(ctype, C.CString(param))
	return C.GoString(charset)
}

func (p *Part) Disposition() string {
	cDisposition := C.gmime_get_content_disposition(p.gmimePart)
	if cDisposition == nil {
		return ""
	}
	//defer C.g_free(C.gpointer(unsafe.Pointer(cDisposition)))
	return C.GoString(cDisposition)
}

// IsText returns true if part's mime is text/*
func (p *Part) IsText() bool {
	return gobool(C.gmime_is_text_part(p.gmimePart))
}

// IsAttachment returns true if part's mime is attachment or inline attachment
func (p *Part) IsAttachment() bool {
	if p.gmimePart == nil {
		return false
	}
	if !gobool(C.gmime_is_part(p.gmimePart)) || gobool(C.gmime_is_multi_part(p.gmimePart)) {
		return false
	}
	if gobool(C.g_mime_part_is_attachment((*C.GMimePart)(unsafe.Pointer(p.gmimePart)))) {
		return true
	}
	if len(p.Filename()) > 0 {
		return true
	}

	return false
}

func (p *Part) IsLegacyAttachment() bool {
	if p.gmimePart == nil {
		return false
	}

	if !gobool(C.gmime_is_part(p.gmimePart)) || gobool(C.gmime_is_multi_part(p.gmimePart)) {
		return false
	}

	if strings.Contains(strings.ToLower(p.Disposition()), "attachment") {
		return true
	}

	contentType := strings.ToLower(p.ContentType())

	if p.parent != nil {
		parentContentType := strings.ToLower(p.parent.ContentType())
		matched, _ := regexp.MatchString(`multipart/(alternative|related|mixed)`, parentContentType)

		// Check if the parent is multipart/alternative and the current part is neither text/plain nor text/html
		if matched && contentType != "text/plain" && contentType != "text/html" {
			return true
		}

		if contentType == parentContentType {
			return false
		}
	}

	if len(p.Filename()) > 0 {
		return true
	}

	return false
}

// Filename retrieves the filename of the part
func (p *Part) Filename() string {
	if p.gmimePart == nil {
		return ""
	}
	if !gobool(C.gmime_is_part(p.gmimePart)) {
		return ""
	}
	ctype := C.g_mime_part_get_filename((*C.GMimePart)(unsafe.Pointer(p.gmimePart)))
	return C.GoString(ctype)
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
	if b == nil {
		return nil
	}
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

// SetHeader sets or replaces specified header
func (p *Part) SetHeader(name string, value string) {
	headers := C.g_mime_object_get_header_list(p.asGMimeObject())
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	cCharset := C.CString("UTF-8")
	defer C.free(unsafe.Pointer(cCharset))

	C.g_mime_header_list_set(headers, cName, cValue, cCharset)
}

// Headers gives you all headers for part
func (p *Part) Headers() textproto.MIMEHeader {
	return nil
}

// ContentID returns the content ID of the attachment if the type is attachment, if not we return empty string
func (p *Part) ContentID() string {
	if p.gmimePart == nil {
		return ""
	}
	if !gobool(C.gmime_is_part(p.gmimePart)) {
		return ""
	}
	cCID := C.g_mime_object_get_content_id(p.gmimePart)
	return C.GoString(cCID)
}

func (p *Part) asGMimeObject() *C.GMimeObject {
	return p.gmimePart
}

// String returns content as a string
func (p *Part) String() string {
	objStr := C.g_mime_object_to_string(p.asGMimeObject(), nil)
	defer C.g_free(C.gpointer(unsafe.Pointer(objStr)))
	return strings.TrimSpace(C.GoString(objStr))
}
