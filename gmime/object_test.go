package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ObjectTestSuite struct {
	suite.Suite
}

func (s *ObjectTestSuite) TestNewObject() {
	contentType := NewContentType("text", "plain")
	assert.Equal(s.T(), contentType.ToString(), "text/plain")
	part := NewObject(contentType)
	assert.NotNil(s.T(), part)
	assert.Equal(s.T(), part.ContentType().ToString(), "text/plain")
	_, ok := part.(Part)
	assert.True(s.T(), ok)
}

func (s *ObjectTestSuite) TestNewObjectWithType() {
	multipart := NewObjectWithType("multipart", "mixed")
	assert.NotNil(s.T(), multipart)
	assert.Equal(s.T(), multipart.ContentType().ToString(), "multipart/mixed")
	_, ok := multipart.(Multipart)
	assert.True(s.T(), ok)
}

func (s *ObjectTestSuite) TestHeader() {
	contentType := NewContentType("text", "plain")
	part := NewObject(contentType)
	part.SetHeader("X-SendGrid-Name", "hola!")
	header, ok := part.Header("X-SendGrid-Name")
	assert.True(s.T(), ok)
	assert.Equal(s.T(), header, "hola!")

	_, ok = part.Header("X-Not-Exists")
	assert.False(s.T(), ok)

	// test multiple headers
	part.SetHeader("X-SendGrid-Name2", "hola2")
	part.SetHeader("X-SendGrid-Name3", "hola3")
	assert.Equal(s.T(), part.Headers(), "Content-Type: text/plain\nX-SendGrid-Name: hola!\nX-SendGrid-Name2: hola2\nX-SendGrid-Name3: hola3\n")
	assert.Equal(s.T(), part.ToString(), "Content-Type: text/plain\nX-SendGrid-Name: hola!\nX-SendGrid-Name2: hola2\nX-SendGrid-Name3: hola3\n\n")
}

func (s *ObjectTestSuite) TestWriteToStream() {
	contentType := NewContentType("text", "plain")
	part := NewObject(contentType)
	part.SetHeader("X-Test", "value")
	stream := NewMemStream()
	part.WriteToStream(stream)
	assert.Equal(s.T(), stream.Length(), 40)
	stream.Close()
}

// run test
func TestObjectTestSuite(t *testing.T) {
	suite.Run(t, new(ObjectTestSuite))
}
