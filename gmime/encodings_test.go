package gmime_test

import (
	"testing"

	"github.com/sendgrid/go-gmime/gmime"
	"github.com/stretchr/testify/assert"
)

const decoded = "Twas brillig, and the slithy toves"
const encoded = "VHdhcyBicmlsbGlnLCBhbmQgdGhlIHNsaXRoeSB0b3Zlcw==\n"

func TestEncodings(t *testing.T) {
	encoder := gmime.NewContentEncoder("base64")
	result := encoder.Flush([]byte(decoded))
	assert.Equal(t, result, []byte(encoded))
}

func TestDecodings(t *testing.T) {
	decoder := gmime.NewContentDecoder("base64")
	result := decoder.Flush([]byte(encoded))
	assert.Equal(t, result, []byte(decoded))
}
