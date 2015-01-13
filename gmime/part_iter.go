package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

type PartIter interface {
	Current() Object
	HasNext() bool
	Next()
}

// XXX Adapted the API to more closely match the Iterator Design Pattern;
// canonical external iterators fit into the idiomatic for-loop usage more naturally.
// TODO Consider using the underlying foreach functions to provide internal iterators.
type aPartIter struct {
	partIter  *C.GMimePartIter
	hasNext   bool
	container Object
}

func NewPartIter(message Message) PartIter {
	rawMessage := message.(rawMessage)
	container, _ := message.MimePart().(Multipart)
	p := &aPartIter{
		partIter:  C.g_mime_part_iter_new((*C.GMimeObject)(unsafe.Pointer(rawMessage.rawMessage()))),
		hasNext:   true,
		container: container,
	}
	runtime.SetFinalizer(p, func(p *aPartIter) { p.free() })
	return p
}

func (p *aPartIter) Current() Object {
	if p.container != nil {
		return p.container
	} else {
		object := C.g_mime_part_iter_get_current(p.partIter)
		if object == nil {
			return nil
		}
		return objectAsSubclass(object)
	}
}

func (p *aPartIter) HasNext() bool {
	return p.hasNext
}

func (p *aPartIter) next() bool {
	if p.container != nil {
		p.container = nil
		return true
	}
	return gobool(C.g_mime_part_iter_next(p.partIter))
}

func (p *aPartIter) Next() {
	p.hasNext = p.next()
}

func (p *aPartIter) free() {
	C.g_mime_part_iter_free(p.partIter)
	p.partIter = nil
}
