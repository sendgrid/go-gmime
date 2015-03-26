package stdio_test

import (
	"bytes"
	"github.com/sendgrid/go-gmime/gmime/stdio"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCioRead(t *testing.T) {
	hello := []byte("Hello, World!")
	rdr := bytes.NewReader(hello)
	wrapped, err := stdio.WrapReader(rdr, false)
	assert.NoError(t, err)
	wf := wrapped.File()
	defer wf.Close()
	siz := len(hello)
	n, buf := wf.Read(1, siz)
	assert.Equal(t, n, siz)
	assert.Equal(t, buf[0:siz], hello)
}

func TestCioWrite(t *testing.T) {
	hello := "Hello, world! I would like to write a number, but printf isn't supported in CGO."
	var buf bytes.Buffer
	wrap, err := stdio.WrapReadWriter(&buf, false)
	assert.NoError(t, err)
	wf := wrap.File()
	defer wf.Close()
	x := wf.Puts(hello)
	assert.NotEqual(t, x, stdio.EOF)
	x = wf.Flush()
	assert.NotEqual(t, x, stdio.EOF)
	assert.Equal(t, string(buf.Bytes()), hello)
}

func TestCioNoClose(t *testing.T) {
	f, err := os.Create("/tmp/stdio_foo.txt")
	assert.NoError(t, err)
	defer f.Close()
	defer os.Remove("/tmp/stdio_foo.txt")
	wrap, err := stdio.WrapWriter(f, false)
	assert.NoError(t, err)
	wf := wrap.File()
	x := wf.Puts("Foo!\n")
	assert.NotEqual(t, x, stdio.EOF)
	wf.Close()
	// FIXME: return value of C.fclose() was checked in orig version
	// assert.NotEqual(t, x, stdio.EOF)
}
