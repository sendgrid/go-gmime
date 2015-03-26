// +build linux
package stdio

/*
#define _GNU_SOURCE
#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <errno.h>
#include <sys/errno.h>
#include "indirect_linux.h"

static void seterr ( int e ) { errno = e; }

static FILE * stdio_setup_cookie(void *cookie, const char *mode, bool r, bool w, bool s) {
    cookie_io_functions_t funcswrap = {
       .read = r ? c_reader : NULL,
       .write = w ? c_writer : NULL,
       .seek =  s ? c_seeker : NULL,
       .close = c_closer,
   };
   return fopencookie(cookie, mode, funcswrap);
}
*/
import "C"

import (
	"errors"
	"io"
	"os"
	"unsafe"
)

//export reader
func reader(cookiePtr unsafe.Pointer, buf *C.char, size C.size_t) C.ssize_t {
	cookie := (*Wrapper)(cookiePtr)
	rdr, ok := cookie.Reader()
	if !ok {
		C.seterr(C.EBADF)
		return -1
	}
	buffer := makeSlice(buf, int(size))
	n, err := rdr.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return C.ssize_t(n)
		}
		C.seterr(C.EIO)
		return -1
	}
	return C.ssize_t(n)
}

//export writer
func writer(cookiePtr unsafe.Pointer, buf *C.char, size C.size_t) C.ssize_t {
	cookie := (*Wrapper)(cookiePtr)
	rdr, ok := cookie.Writer()
	if !ok {
		C.seterr(C.EBADF)
		return -1
	}
	buffer := makeSlice(buf, int(size))
	n, err := rdr.Write(buffer)
	if err != nil {
		if err == io.EOF {
			return C.ssize_t(n)
		}
		C.seterr(C.EIO)
		return -1
	}
	return C.ssize_t(n)
}

//export closer
func closer(cookiePtr unsafe.Pointer) C.int {
	cookie := (*Wrapper)(cookiePtr)
	cls, ok := cookie.Closer()
	var rc C.int = 0
	if ok {
		if err := cls.Close(); err != nil {
			C.seterr(C.EIO)
			rc = -1
		}
	}
	cookie.closed = true
	return rc
}

//export seeker
func seeker(cookiePtr unsafe.Pointer, position *C.off64_t, whence C.int) C.int {
	cookie := (*Wrapper)(cookiePtr)
	skr, ok := cookie.Seeker()
	if !ok {
		C.seterr(C.EBADF)
		return -1
	}
	var w int
	// Not sure if C.SEEK_* matches os.SEEK_* in all cases.
	switch whence {
	case C.SEEK_SET:
		w = os.SEEK_SET
	case C.SEEK_CUR:
		w = os.SEEK_CUR
	case C.SEEK_END:
		w = os.SEEK_END
	default:
		C.seterr(C.EINVAL)
		return -1
	}
	ret, err := skr.Seek(int64(*position), w)
	if err != nil {
		C.seterr(C.EINVAL)
		return -1
	}
	*position = C.off64_t(ret)
	return 0
}

func wrapReadWriter(cookie *Wrapper) (*Wrapper, error) {
	_, skr := cookie.val.(io.Seeker)
	_, rdr := cookie.val.(io.Reader)
	_, wtr := cookie.val.(io.Writer)
	mode := "w+"
	if !rdr {
		mode = "w"
	}
	if !wtr {
		mode = "r"
	}
	cmode := C.CString(mode)
	defer C.free(unsafe.Pointer(cmode))
	f, err := C.stdio_setup_cookie(unsafe.Pointer(cookie), cmode, C._Bool(rdr), C._Bool(wtr), C._Bool(skr))
	if f == nil {
		return nil, errors.New("fopencookie")
	}
	nf := newFile(f, mode)
	cookie.f = nf
	return cookie, err
}
