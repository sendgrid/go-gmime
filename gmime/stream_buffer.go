package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <string.h>
#include <gmime/gmime.h>
*/
import "C"
import (
	"unsafe"
)

const (
	CACHE_READ  = C.GMIME_STREAM_BUFFER_CACHE_READ
	BLOCK_READ  = C.GMIME_STREAM_BUFFER_BLOCK_READ
	BLOCK_WRITE = C.GMIME_STREAM_BUFFER_BLOCK_WRITE
)

type BufferedStream interface {
	Stream
	Gets(char *[]byte, max int64) int64
	ReadLn(buffer *[]byte)
}

type aBufferedStream struct {
	*aStream
}

func CastBufferedStream(cs *C.GMimeStream) *aBufferedStream {
	return &aBufferedStream{CastStream(cs)}
}

func NewBufferedStream(source Stream, mode C.GMimeStreamBufferMode) BufferedStream {
	rawSource := source.(rawStream)
	cStream := C.g_mime_stream_buffer_new(rawSource.rawStream(), mode)
	defer unref(C.gpointer(cStream))
	return CastBufferedStream(cStream)
}

func (s *aBufferedStream) Gets(char *[]byte, max int64) int64 {
	cLength := C.g_mime_stream_buffer_gets(s.rawStream(), (*C.char)(unsafe.Pointer(&char)), (C.size_t)(max))

	return int64(cLength)
}

func (s *aBufferedStream) ReadLn(buffer *[]byte) {
	C.g_mime_stream_buffer_readln(s.rawStream(), (*C.GByteArray)(unsafe.Pointer(&buffer)))
}
