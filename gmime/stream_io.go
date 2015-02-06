package gmime

import (
	"github.com/sendgrid/go-gmime/gmime/cio"
	"io"
)

type ioStream struct {
	*aFileStream
	wrapper *cio.Wrapper
}

func newIOStream(w *cio.Wrapper, err error) FileStream {
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
	return newIOStream(cio.WrapReader(r, doClose))
}

func NewWriterStream(w io.Writer, doClose bool) FileStream {
	return newIOStream(cio.WrapWriter(w, doClose))
}

func NewReadWriterStream(rw io.ReadWriter, doClose bool) FileStream {
	return newIOStream(cio.WrapReadWriter(rw, doClose))
}
