package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataWrapper(t *testing.T) {
	dw := NewDataWrapper()
	assert.Equal(t, dw.Encoding(), "")
}

func TestNewDataWrapperWithStream(t *testing.T) {
	loop := 1

	for i := 0; i < loop; i++ {
		raw := "foo=bar"
		escaped := "foo\xbar"

		instream := NewMemStream()
		instream.Length()
		instream.WriteString(raw)
		encoding := "quoted-printable"

		wrapper := NewDataWrapperWithStream(instream, encoding)
		assert.Equal(t, wrapper.Encoding(), encoding)

		outstream := NewMemStream()
		wrapper.WriteToStream(outstream)
		assert.Equal(t, string(outstream.Bytes()), escaped)
	}
}

func TestDataWrapperStream(t *testing.T) {
	stream := NewMemStreamWithBuffer("hola")
	encoding := "gzip"
	wrapper := NewDataWrapperWithStream(stream, encoding)
	assert.Equal(t, int64(4), wrapper.Stream().Length())
}
