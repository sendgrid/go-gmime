package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"

import (
	"unsafe"
)

type Filter interface {
	Copy() Filter
	Reset()
}

// FIXME: Merge with Filter when we implement this
type FilterFullLowlevelInterface interface {
	Filter(inbuf []byte, prespace int) (outbuf []byte, outlen int, outprespace int)
	Complete(inbuf []byte, prespace int) (outbuf []byte, outlen int, outprespace int)
}

type gmimeFilter interface {
	Janitor
	rawFilter() *C.GMimeFilter
}

type aFilter struct {
	*PointerMixin
}

func castFilter(ptr *C.GMimeFilter) *aFilter {
	return &aFilter{CastPointer(C.gpointer(ptr))}
}

func (f *aFilter) Copy() *aFilter {
	return castFilter(C.g_mime_filter_copy(f.rawFilter()))
}

func (f *aFilter) Reset() {
	C.g_mime_filter_reset(f.rawFilter())
}

func (f *aFilter) rawFilter() *C.GMimeFilter {
	return (*C.GMimeFilter)(f.pointer())
}

/// Specializing

func NewBasicFilter(encoding string, encode bool) *aFilter {
	enc := goGMimeString2Encoding(encoding)
	f := C.g_mime_filter_basic_new(enc, gbool(encode))
	return castFilter(f)
}

type FilterBest interface {
	Charset() string
	Encoding() string
}

type filterBest interface {
	gmimeFilter
}

type aBestFilter struct {
	*aFilter
}

// type GMimeFilterBestFlags int
const (
	GMIME_FILTER_BEST_CHARSET  = (1 << 0)
	GMIME_FILTER_BEST_ENCODING = (1 << 1)
)

func NewBestFilter(flags int) *aBestFilter {
	f := castFilter((*C.GMimeFilter)(C.g_mime_filter_best_new(C.GMimeFilterBestFlags(flags))))
	return &aBestFilter{f}
}

func (f *aBestFilter) Charset() string {
	s := C.g_mime_filter_best_charset((*C.GMimeFilterBest)(f.pointer()))
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// type GMimeEncodingConstraint int
const (
	GMIME_ENCODING_CONSTRAINT_7BIT   = C.GMIME_ENCODING_CONSTRAINT_7BIT
	GMIME_ENCODING_CONSTRAINT_8BIT   = C.GMIME_ENCODING_CONSTRAINT_8BIT
	GMIME_ENCODING_CONSTRAINT_BINARY = C.GMIME_ENCODING_CONSTRAINT_BINARY
)

func (f *aBestFilter) Encoding(constraint int) string {
	e := C.g_mime_filter_best_encoding((*C.GMimeFilterBest)(f.pointer()), C.GMimeEncodingConstraint(constraint))
	return goGMimeEncoding2String(e)
}
