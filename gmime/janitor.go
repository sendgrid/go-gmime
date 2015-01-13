package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>

*/
import "C"

import (
	"runtime"
)

// Janitor to clean up wasteful resources
type Janitor interface {
	pointer() C.gpointer
	finalize()
}

func AssignJanitor(self Janitor) {
	ref(self.pointer())
	runtime.SetFinalizer(self, func(j Janitor) {
		j.finalize()
	})
}

type PointerMixin struct {
	ptr C.gpointer
}

func CastPointer(p C.gpointer) *PointerMixin {
	pp := &PointerMixin{ptr: p}
	AssignJanitor(pp)
	return pp
}

// for newly allocated only
func NewPointer(p C.gpointer) *PointerMixin {
	defer unref(p)
	return CastPointer(p)
}

func (self *PointerMixin) finalize() {
	unref(self.ptr)
}

func (self *PointerMixin) pointer() C.gpointer {
	return self.ptr
}
