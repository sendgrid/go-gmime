package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>

static GMimeEncoding * alloc_gmime_encoding() {
    return g_new0(GMimeEncoding, 1);
}
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type ContentEncodingState interface {
	Outlen(inlen int) (outlen int)
	Step(in []byte) (out []byte)
	Flush(in []byte) (out []byte)
}

type contentEncoding struct {
	Pointer *C.GMimeEncoding
}

func allocateContentEncoder() *contentEncoding {
	ptr := C.alloc_gmime_encoding()
	obj := &contentEncoding{
		Pointer: ptr,
	}
	runtime.SetFinalizer(obj, func(o *contentEncoding) {
		C.g_free((C.gpointer)(unsafe.Pointer(o.Pointer)))
	})
	return obj
}

func goGMimeString2Encoding(encoding string) C.GMimeContentEncoding {
	cEncoding := C.CString(encoding)
	defer C.free(unsafe.Pointer(cEncoding))
	return C.g_mime_content_encoding_from_string(cEncoding)
}

func goGMimeEncoding2String(enc C.GMimeContentEncoding) string {
	return C.GoString(C.g_mime_content_encoding_to_string(enc))
}

func NewContentEncoder(encoding string) ContentEncodingState {
	e := allocateContentEncoder()
	C.g_mime_encoding_init_encode(e.Pointer, goGMimeString2Encoding(encoding))
	return ContentEncodingState(e)
}

func NewContentDecoder(encoding string) ContentEncodingState {
	e := allocateContentEncoder()
	C.g_mime_encoding_init_decode(e.Pointer, goGMimeString2Encoding(encoding))
	return ContentEncodingState(e)
}

func (e *contentEncoding) Outlen(inlen int) int {
	return int(C.g_mime_encoding_outlen(e.Pointer, C.size_t(inlen)))
}

func (e *contentEncoding) Step(in []byte) []byte {
	l := len(in)
	outlen := e.Outlen(l)
	inbuf := (*C.char)(unsafe.Pointer(&in[0]))
	out := make([]byte, outlen)
	outbuf := (*C.char)(unsafe.Pointer(&out[0]))
	rlen := int(C.g_mime_encoding_step(e.Pointer, inbuf, C.size_t(l), outbuf))
	return out[0:rlen]
}

func (e *contentEncoding) Flush(in []byte) []byte {
	l := len(in)
	outlen := e.Outlen(l)
	inbuf := (*C.char)(unsafe.Pointer(&in[0]))
	out := make([]byte, outlen)
	outbuf := (*C.char)(unsafe.Pointer(&out[0]))
	rlen := int(C.g_mime_encoding_flush(e.Pointer, inbuf, C.size_t(l), outbuf))
	return out[0:rlen]
}
