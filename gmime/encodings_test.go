package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodings(t *testing.T) {
	sevenBit := NewContentEncodingFromString("7bit")
	assert.Equal(t, sevenBit.ToString(), "7bit")

	eightBit := NewContentEncodingFromString("8bit")
	assert.Equal(t, eightBit.ToString(), "8bit")

	binary := NewContentEncodingFromString("binary")
	assert.Equal(t, binary.ToString(), "binary")

	base64 := NewContentEncodingFromString("base64")
	assert.Equal(t, base64.ToString(), "base64")

	qp := NewContentEncodingFromString("quoted-printable")
	assert.Equal(t, qp.ToString(), "quoted-printable")

	// an extension, so prefixed with x-
	uuencode := NewContentEncodingFromString("uuencode")
	assert.Equal(t, uuencode.ToString(), "x-uuencode")

	_default := NewContentEncodingFromString("")
	assert.Equal(t, _default.ToString(), "")

	invalid := NewContentEncodingFromString("invalid")
	assert.Equal(t, invalid.ToString(), _default.ToString())
}
