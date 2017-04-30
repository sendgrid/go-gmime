package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PartTestSuite struct {
	suite.Suite
}

func (s *PartTestSuite) TestNewPart() {
	contentType := NewContentType("text", "plain")
	part := NewPart()
	assert.NotNil(s.T(), part)
	part.SetContentType(contentType)
	assert.Equal(s.T(), part.ContentType().ToString(), "text/plain")

	instream := NewMemStreamWithBuffer("text")
	outstream := NewMemStream() // NewMemStreamWithByteArray()
	wrapper := NewDataWrapperWithStream(instream, "7bit")
	part.SetContentObject(wrapper)

	part.ContentObject().WriteToStream(outstream)
	assert.Equal(s.T(), (string)(outstream.Bytes()), "text")
	assert.Equal(s.T(), part.ContentObject().Encoding(), "7bit")
}

func (s *PartTestSuite) TestNewPartWithType() {
	part := NewPartWithType("text", "html")
	assert.NotNil(s.T(), part)
	assert.Equal(s.T(), part.ContentType().ToString(), "text/html")

	instream := NewMemStreamWithBuffer("<html></html>")
	outstream := NewMemStream() // NewMemStreamWithByteArray()
	wrapper := NewDataWrapperWithStream(instream, "8bit")
	part.SetContentObject(wrapper)

	part.ContentObject().WriteToStream(outstream)
	assert.Equal(s.T(), (string)(outstream.Bytes()), "<html></html>")
	assert.Equal(s.T(), part.ContentObject().Encoding(), "8bit")
}

func (s *PartTestSuite) TestContentObject() {
	stream := NewMemStreamWithBuffer("hola")
	dw := NewDataWrapperWithStream(stream, "gzip")
	part := NewPart()
	part.SetContentObject(dw)
	dwFromContent := part.ContentObject()
	length := dwFromContent.Stream().Length()
	assert.Equal(s.T(), int64(4), length)
}

func (s *PartTestSuite) TestFilename() {
	part := NewPart()
	assert.Equal(s.T(), part.Filename(), "")
}

func (s *PartTestSuite) TestDescription() {
	part := NewPart()
	assert.Equal(s.T(), part.Description(), "")
}

func (s *PartTestSuite) TestContentLocation() {
	part := NewPart()
	assert.Equal(s.T(), part.ContentLocation(), "")
}

func (s *PartTestSuite) TestContentEncoding() {
	part := NewPart()
	part.SetContentEncoding("gzip")
	assert.Equal(s.T(), part.ContentEncoding(), "") // FIXME: ???
}

func TestPartTestSuite(t *testing.T) {
	suite.Run(t, new(PartTestSuite))
}
