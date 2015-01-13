package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"

import (
	"os"
	"unsafe"
)

type FileStream interface {
	Stream
	Owner() bool
	SetOwner(bool)
}

type aFileStream struct {
	*aStream
}

func CastFileStream(cfs *C.GMimeStreamFile) *aFileStream {
	fs := (*C.GMimeStream)(unsafe.Pointer(cfs))
	s := CastStream(fs)
	return &aFileStream{s}
}

func NewFileStream(f *os.File) FileStream {
	return NewFileStreamWithMode(f, "a")
}

func NewFileStreamWithMode(f *os.File, mode string) FileStream {
	cMode := C.CString(mode)
	defer C.free(unsafe.Pointer(cMode))
	cFile := C.fdopen(C.int(f.Fd()), cMode)
	s := C.g_mime_stream_file_new(cFile)
	fileStream := (*C.GMimeStreamFile)(unsafe.Pointer(s))
	defer unref(C.gpointer(fileStream))
	return CastFileStream(fileStream)
}

func NewFileStreamWithBounds(f os.File, start int64, end int64) FileStream {
	mode := C.CString("r")
	defer C.free(unsafe.Pointer(mode))
	cFile := C.fdopen(C.int(f.Fd()), mode)
	sBound := C.g_mime_stream_file_new_with_bounds(cFile, (C.gint64)(start), (C.gint64)(end))
	fileStream := (*C.GMimeStreamFile)(unsafe.Pointer(sBound))
	defer unref(C.gpointer(fileStream))
	return CastFileStream(fileStream)
}

func (f *aFileStream) rawFileStream() *C.GMimeStreamFile {
	return (*C.GMimeStreamFile)(f.pointer())
}

func (f *aFileStream) Owner() bool {
	result := C.g_mime_stream_file_get_owner(f.rawFileStream())
	return gobool(result)
}

func (f *aFileStream) SetOwner(owner bool) {
	C.g_mime_stream_file_set_owner(f.rawFileStream(), gbool(owner))
}
