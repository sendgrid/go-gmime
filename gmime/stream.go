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

type Stream interface {
	Janitor
	Read(int64) (int64, []byte)
	Write([]byte) int64
	WriteBlock([]byte, int64) int64
	Flush() int
	Close() int
	Eos() bool
	Reset() int
	Seek(int64, int) (int64, error)
	Tell() int64
	Length() int64
	SubStream(int64, int64) Stream
	SetBounds(int64, int64)
	WriteString(string) int64
	WriteToStream(Stream) int64
	//	Raw() *C.GMimeStream
}

type aStream struct {
	*PointerMixin
}

type rawStream interface {
	Stream
	rawStream() *C.GMimeStream
}

func CastStream(cs *C.GMimeStream) *aStream {
	return &aStream{CastPointer(C.gpointer(cs))}
}

func (s *aStream) Read(length int64) (int64, []byte) {
	cBuf := make([]byte, length)
	cChar := (*C.char)(unsafe.Pointer(&cBuf[0]))

	cLength := C.g_mime_stream_read(s.rawStream(), cChar, (C.size_t)(length))
	if cLength <= 0 {
		return int64(cLength), nil
	}
	return int64(cLength), cBuf[:cLength]
}

func (s *aStream) Write(buf []byte) int64 {
	return s.WriteBlock(buf, int64(len(buf)))
}

func (s *aStream) WriteBlock(buf []byte, length int64) int64 {
	cLength := C.g_mime_stream_write(s.rawStream(), (*C.char)(unsafe.Pointer(&buf[0])), (C.size_t)(length))
	return int64(cLength)
}

func (s *aStream) Flush() int {
	return int(C.g_mime_stream_flush(s.rawStream()))
}

func (s *aStream) Close() int {
	return int(C.g_mime_stream_close(s.rawStream()))
}

func (s *aStream) Eos() bool {
	ret := C.g_mime_stream_eos(s.rawStream())
	return gobool(ret)
}

func (s *aStream) Reset() int {
	return int(C.g_mime_stream_reset(s.rawStream()))
}

func (s *aStream) Seek(offset int64, whence int) (int64, error) {
	length, errno := C.g_mime_stream_seek(s.rawStream(), (C.gint64)(offset), (C.GMimeSeekWhence)(whence))
	return int64(length), error(errno)
}

func (s *aStream) Tell() int64 {
	return int64(C.g_mime_stream_tell(s.rawStream()))
}

func (s *aStream) Length() int64 {
	if s.rawStream() != nil {
		return int64(C.g_mime_stream_length(s.rawStream()))
	}
	return int64(0)
}

func (s *aStream) SubStream(start int64, end int64) Stream {
	cStream := C.g_mime_stream_substream(s.rawStream(), (C.gint64)(start), (C.gint64)(end))
	defer unref(C.gpointer(cStream))

	subStream := CastStream(cStream)
	return subStream
}

func (s *aStream) SetBounds(start int64, end int64) {
	C.g_mime_stream_set_bounds(s.rawStream(), (C.gint64)(start), (C.gint64)(end))
}

func (s *aStream) WriteString(str string) int64 {
	var cStr *C.char = C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	res := C.g_mime_stream_write_string(s.rawStream(), cStr)
	return int64(res)
}

func (s *aStream) WriteToStream(dest Stream) int64 {
	rawDest := dest.(rawStream)
	res := C.g_mime_stream_write_to_stream(s.rawStream(), rawDest.rawStream())
	return int64(res)
}

func (s *aStream) rawStream() *C.GMimeStream {
	return (*C.GMimeStream)(s.pointer())
}
