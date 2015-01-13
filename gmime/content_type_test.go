package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContentType(t *testing.T) {
	loop := 1
	for i := 0; i < loop; i++ {
		ct := NewContentType("text", "plain")
		assert.Equal(t, ct.ToString(), "text/plain")
		assert.Equal(t, ct.MediaType(), "text")
		assert.Equal(t, ct.MediaSubtype(), "plain")

		defaultType := NewContentType("", "")
		assert.Equal(t, defaultType.ToString(), "application/octet-stream")
		assert.Equal(t, defaultType.MediaType(), "application")
		assert.Equal(t, defaultType.MediaSubtype(), "octet-stream")

		customType := NewContentType("made", "up")
		assert.Equal(t, customType.ToString(), "made/up")
		assert.Equal(t, customType.MediaType(), "made")
		assert.Equal(t, customType.MediaSubtype(), "up")

		multipartType := NewContentType("multipart", "mixed")
		boundary := "--foo==bar++foo--"
		multipartType.SetParameter("boundary", boundary)
		assert.Equal(t, multipartType.ToString(), "multipart/mixed")
		assert.Equal(t, multipartType.MediaType(), "multipart")
		assert.Equal(t, multipartType.MediaSubtype(), "mixed")
		assert.Equal(t, multipartType.Parameter("boundary"), boundary)
	}
}

func TestNewContentTypeFromString(t *testing.T) {
	loop := 1
	for i := 0; i < loop; i++ {
		ct := NewContentTypeFromString("text/plain")
		assert.Equal(t, ct.ToString(), "text/plain")
		assert.Equal(t, ct.MediaType(), "text")
		assert.Equal(t, ct.MediaSubtype(), "plain")

		defaultType := NewContentTypeFromString("")
		assert.Equal(t, defaultType.ToString(), "application/octet-stream")
		assert.Equal(t, defaultType.MediaType(), "application")
		assert.Equal(t, defaultType.MediaSubtype(), "octet-stream")

		customType := NewContentTypeFromString("made/up")
		assert.Equal(t, customType.ToString(), "made/up")
		assert.Equal(t, customType.MediaType(), "made")
		assert.Equal(t, customType.MediaSubtype(), "up")

		multipartType := NewContentTypeFromString("multipart/mixed")
		boundary := "--foo==bar++foo--"
		multipartType.SetParameter("boundary", boundary)
		assert.Equal(t, multipartType.ToString(), "multipart/mixed")
		assert.Equal(t, multipartType.MediaType(), "multipart")
		assert.Equal(t, multipartType.MediaSubtype(), "mixed")
		assert.Equal(t, multipartType.Parameter("boundary"), boundary)
	}
}

func TestContentTypeToString(t *testing.T) {
	ct := NewContentTypeFromString("text/plain")
	assert.Equal(t, ct.ToString(), "text/plain")
}

func TestContentTypeSetParameter(t *testing.T) {
	ct := NewContentTypeFromString("application/json")
	ct.SetParameter("hola", "hi!")
	assert.Equal(t, ct.Parameter("hola"), "hi!")
}

/*
FIXME: rewrite test
func TestContentTypeParams(t *testing.T) {
	ct := NewContentTypeFromString("application/json")
	ct.SetParameter("hola", "hi!")
	assert.Equal(t, ct.Params().Name(), "hola")
	assert.Equal(t, ct.Params().Value(), "hi!")

	// set multiple params:
	ct.SetParameter("key1", "val1")
	ct.SetParameter("key2", "val2")
	params := ct.Params()
	params = params.Next()
	assert.Equal(t, params.Name(), "key1")
	assert.Equal(t, params.Value(), "val1")
	params = params.Next()
	assert.Equal(t, params.Name(), "key2")
	assert.Equal(t, params.Value(), "val2")
}
*/

func TestContentTypeMediaType(t *testing.T) {
	ct := NewContentTypeFromString("application/json")
	assert.Equal(t, ct.MediaType(), "application")
	assert.Equal(t, ct.MediaSubtype(), "json")
}
