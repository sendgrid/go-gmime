package gmime

import (
	"github.com/sendgrid/go-gmime/gmime/stdio"
	"io"
)

type ioStream struct {
	*aFileStream
	wrapper *stdio.Wrapper
}

func newIOStream(w *stdio.Wrapper, err error) FileStream {
	if w == nil {
		return nil
	}
	file := w.File()
	fs := NewFileStreamWithMode(file.Pointer(), file.Mode())
	return &ioStream{
		aFileStream: fs.(*aFileStream),
		wrapper:     w,
	}
}

func NewReaderStream(r io.Reader, doClose bool) FileStream {
	return newIOStream(stdio.WrapReader(r, doClose))
}

func NewWriterStream(w io.Writer, doClose bool) FileStream {
	return newIOStream(stdio.WrapWriter(w, doClose))
}

func NewReadWriterStream(rw io.ReadWriter, doClose bool) FileStream {
	return newIOStream(stdio.WrapReadWriter(rw, doClose))
}
