package gmime

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MultipartTestSuite struct {
	suite.Suite
	part Multipart
}

func (s *MultipartTestSuite) SetupTest() {
	s.part = NewMultipart()
}

func (s *MultipartTestSuite) TearDownTest() {
}

func (s *MultipartTestSuite) TestNewMultipart() {
	assert.NotNil(s.T(), s.part)
}

func (s *MultipartTestSuite) TestPart() {
	part := NewPart()
	wrapper := NewDataWrapper()
	part.SetContentObject(wrapper)
	s.part.AddPart(part)

	allParts := s.part.GetPart(0)
	_, ok := allParts.(Part)
	assert.True(s.T(), ok)
}

func (s *MultipartTestSuite) TestContentType() {
	typeString := "multipart/mixed"
	boundary := "--foo++bar==foo--"
	contentType := NewContentTypeFromString(typeString)
	contentType.SetParameter("boundary", boundary)
	s.part.SetContentType(contentType)

	assert.Equal(s.T(), s.part.ContentType().ToString(), typeString)
	assert.Equal(s.T(), s.part.ContentType().Parameter("boundary"), boundary)
}

func (s *MultipartTestSuite) TestPartGC() {
	part := NewPart()
	wrapper := NewDataWrapper()
	part.SetContentObject(wrapper)
	s.part.AddPart(part)

	allParts := s.part.GetPart(0)
	allParts2 := s.part.GetPart(0)

	_, ok := allParts.(Part)
	_, ok = allParts2.(Part)

	allParts = nil
	part = nil
	wrapper = nil
	debug.FreeOSMemory()

	// Try to access allParts2.
	// This should not fail under valgrind
	typeString := "multipart/mixed"
	boundary := "--foo++bar==foo--"
	contentType := NewContentTypeFromString(typeString)
	contentType.SetParameter("boundary", boundary)
	allParts2.SetContentType(contentType)

	assert.True(s.T(), ok)
}

// run test
func TestMultipartTestSuite(t *testing.T) {
	suite.Run(t, new(MultipartTestSuite))
}
