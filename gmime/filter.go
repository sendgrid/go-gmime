package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"


type Filter interface {
    Copy() Filter
    Reset()
}

// FIXME: Merge with Filter when we implement this
type FilterFullLowlevelInterface interface {
    Filter(inbuf []byte, prespace int) (outbuf []byte, outlen int, outprespace int)
    Complete(inbuf []byte, prespace int) (outbuf []byte, outlen int, outprespace int)
}

type gmimeFilter_ interface {
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
