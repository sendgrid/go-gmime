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
	fc := C.g_mime_filter_copy(f.rawFilter())
	defer unref(C.gpointer(fc))
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
	defer unref(C.gpointer(f))
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
	defer unref(C.gpointer(f))
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

func NewCharsetFilter(source string, target string) *aFilter {
	cSrc := C.CString(source)
	defer C.free(unsafe.Pointer(cSrc))
	cTgt := C.CString(target)
	defer C.free(unsafe.Pointer(cTgt))
	f := C.g_mime_filter_charset_new(cSrc, cTgt)
	return castFilter(f)
}

func NewCRLFFilter(encode bool, dots bool) *aFilter {
	f := C.g_mime_filter_crlf_new(gbool(encode), gbool(dots))
	defer unref(C.gpointer(f))
	return castFilter(f)
}

const (
	GMIME_FILTER_ENRICHED_IS_RICHTEXT int32 = (1 << 0)
)

func NewEnrichedFilter(flags int32) *aFilter {
	f := C.g_mime_filter_enriched_new(C.guint32(flags))
	return castFilter(f)
}

const (
	GMIME_FILTER_FROM_MODE_DEFAULT = 0
	GMIME_FILTER_FROM_MODE_ESCAPE  = 0
	GMIME_FILTER_FROM_MODE_ARMOR   = 0
)

func MewFromFilter(flags int) *aFilter {
	f := C.g_mime_filter_from_new(C.GMimeFilterFromMode(flags))
	return castFilter(f)
}

const (
	GMIME_FILTER_GZIP_MODE_ZIP   = C.GMIME_FILTER_GZIP_MODE_ZIP
	GMIME_FILTER_GZIP_MODE_UNZIP = C.GMIME_FILTER_GZIP_MODE_UNZIP
)

func NewGZipFilter(mode int, level int) *aFilter {
	f := C.g_mime_filter_gzip_new(C.GMimeFilterGZipMode(mode), C.int(level))
	return castFilter(f)
}

const (
	GMIME_FILTER_HTML_PRE               uint32 = (1 << 0)
	GMIME_FILTER_HTML_CONVERT_NL               = (1 << 1)
	GMIME_FILTER_HTML_CONVERT_SPACES           = (1 << 2)
	GMIME_FILTER_HTML_CONVERT_URLS             = (1 << 3)
	GMIME_FILTER_HTML_MARK_CITATION            = (1 << 4)
	GMIME_FILTER_HTML_CONVERT_ADDRESSES        = (1 << 5)
	GMIME_FILTER_HTML_ESCAPE_8BIT              = (1 << 6)
	GMIME_FILTER_HTML_CITE                     = (1 << 7)
)

func NewHTMLFilter(flags uint32, colour int32) *aFilter {
	f := C.g_mime_filter_html_new(C.guint32(flags), C.guint32(colour))
	defer unref(C.gpointer(f))
	return castFilter(f)
}

func NewStripFilter() *aFilter {
	f := C.g_mime_filter_strip_new()
	defer unref(C.gpointer(f))
	return castFilter(f)
}

type MD5Filter interface {
	Filter
	GetDigest() []byte
}

type aFilterMD5 struct {
	*aFilter
}

func NewMD5Filter() *aFilterMD5 {
	f := castFilter(C.g_mime_filter_md5_new())
	defer unref(C.gpointer(f))
	return &aFilterMD5{f}
}

func (m *aFilterMD5) GetDigest() []byte {
	b := make([]byte, 16)
	bp := &b[0]
	mp := (*C.GMimeFilterMd5)(unsafe.Pointer(m.pointer()))
	C.g_mime_filter_md5_get_digest(mp, (*C.uchar)(unsafe.Pointer(bp)))
	return b
}

func NewYEncFilter(encode bool) *aFilter {
	f := C.g_mime_filter_yenc_new(gbool(encode))
	defer unref(C.gpointer(f))
	return castFilter(f)
}
