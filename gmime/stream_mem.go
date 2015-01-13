package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <string.h>
#include <gmime/gmime.h>
static GMimeStream * gmime_mem_stream_to_stream(GMimeStreamMem *obj) {
	return GMIME_STREAM(obj);
}
*/
import "C"

import (
	"unsafe"
)

type MemStream interface {
	Stream
	Bytes() []byte
}

type aMemStream struct {
	*aStream
}

func CastMemStream(cms *C.GMimeStreamMem) *aMemStream {
	s := CastStream((*C.GMimeStream)(unsafe.Pointer(cms)))
	return &aMemStream{s}
}

func NewMemStream() MemStream {
	cStream := C.g_mime_stream_mem_new()
	cMemStream := (*C.GMimeStreamMem)(unsafe.Pointer(cStream))
	defer unref(C.gpointer(cMemStream))
	return CastMemStream(cMemStream)
}

func NewMemStreamWithBuffer(str string) MemStream {
	var cBuffer *C.char = C.CString(str)
	defer C.free(unsafe.Pointer(cBuffer))

	strLen := len(str)
	cStream := C.g_mime_stream_mem_new_with_buffer(cBuffer, (C.size_t)(strLen))
	cMemStream := (*C.GMimeStreamMem)(unsafe.Pointer(cStream))
	defer unref(C.gpointer(cMemStream))
	return CastMemStream(cMemStream)
}

func (m *aMemStream) Bytes() []byte {
	ptr := (*C.GMimeStreamMem)(m.pointer())
	if m.Length() > 1 {
		bArray := C.g_mime_stream_mem_get_byte_array(ptr)
		length := bArray.len
		cBuf := make([]byte, length)
		cChar := unsafe.Pointer(unsafe.Pointer(&cBuf[0]))
		//        defer C.g_byte_array_unref(bArray)
		C.memcpy(cChar, unsafe.Pointer(bArray.data), C.size_t(length))
		return cBuf
	}

	return nil
}
