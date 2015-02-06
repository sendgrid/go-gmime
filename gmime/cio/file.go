package cio

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

const (
	EOF int = C.EOF
)

type File struct {
	f    *C.FILE
	mode string
}

// private constructor
func newFile(f *C.FILE, mode string) *File {
	return &File{
		f:    f,
		mode: mode,
	}
}

func (f File) Pointer() unsafe.Pointer {
	return unsafe.Pointer(f.f)
}

func (f File) Mode() string {
	return f.mode
}

func (f File) Error() int {
	return int(C.ferror(f.f))
}

func (f File) Close() {
	C.fclose(f.f)
}

func (f File) Flush() int {
	return int(C.fflush(f.f))
}

func (f File) Read(size int, n int) (int, []byte) {
	b := make([]byte, size*n)
	bp := unsafe.Pointer(&b[0])
	r := C.fread(bp, C.size_t(size), C.size_t(n), f.f)
	return int(r), b
}

func (f File) Write(b []byte, size int, n int) int {
	return int(C.fwrite(unsafe.Pointer(&b[0]), C.size_t(size), C.size_t(n), f.f))
}

func (f File) WriteBytes(b []byte) int {
	return f.Write(b, len(b), 1)
}

func (f File) Puts(s string) int {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return int(C.fputs(cs, f.f))
}
