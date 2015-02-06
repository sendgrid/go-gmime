package cio

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"io"
	"runtime"
)

type Wrapper struct {
	val     interface{}
	f       *File
	mode    string
	doClose bool
	closed  bool
}

func newWrapper(val interface{}, doClose bool) *Wrapper {
	w := &Wrapper{
		val:     val,
		doClose: doClose,
		closed:  false,
		f:       nil,
	}
	runtime.SetFinalizer(w, func(ww *Wrapper) {
		if !ww.closed && ww.f != nil {
			ww.f.Close()
		}
	})
	return w
}

func (w Wrapper) File() *File {
	return w.f
}

func (c Wrapper) Closer() (io.Closer, bool) {
	x, ok := c.val.(io.Closer)
	return x, c.doClose && ok
}
func (c Wrapper) Writer() (io.Writer, bool) {
	x, ok := c.val.(io.Writer)
	return x, ok
}
func (c Wrapper) Seeker() (io.Seeker, bool) {
	x, ok := c.val.(io.Seeker)
	return x, ok
}
func (c Wrapper) Reader() (io.Reader, bool) {
	x, ok := c.val.(io.Reader)
	return x, ok
}
