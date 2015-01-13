package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MessagePartTestSuite struct {
	suite.Suite
	part    MessagePart
	subtype string
}

func (s *MessagePartTestSuite) SetupTest() {
}

func (s *MessagePartTestSuite) TearDownTest() {
}

func (s *MessagePartTestSuite) TestNewMessagePart() {
	subtype := "rfc822"
	part := NewMessagePart(subtype)
	assert.NotNil(s.T(), part)
	rawMessagePart := part.(rawMessagePart)
	assert.NotNil(s.T(), rawMessagePart.rawMessagePart())
}

func (s *MessagePartTestSuite) TestNewMessagePartWithMessage() {
	subtype := "rfc822"
	message := NewMessage()
	part := NewMessagePartWithMessage(subtype, message)
	assert.NotNil(s.T(), part)
	assert.Equal(s.T(), subtype, part.ContentType().MediaSubtype())
	assert.NotNil(s.T(), part.Message())
	assert.Equal(s.T(), message, part.Message())
}

func (s *MessagePartTestSuite) TestMessage() {
	subtype := "rfc822"
	message := NewMessage()
	part := NewMessagePartWithMessage(subtype, message)
	part.SetMessage(message)
	assert.NotNil(s.T(), part)
	assert.Equal(s.T(), subtype, part.ContentType().MediaSubtype())
	assert.NotNil(s.T(), part.Message())
	assert.Equal(s.T(), message, part.Message())
}

func (s *MessagePartTestSuite) TestContentType() {
	subtype := "rfc822"
	part := NewMessagePart(subtype)
	typeString := "message/rfc822"
	contentType := NewContentTypeFromString(typeString)
	part.SetContentType(contentType)

	assert.Equal(s.T(), part.ContentType().ToString(), typeString)
}

// run test
func TestMessagePartTestSuite(t *testing.T) {
	suite.Run(t, new(MessagePartTestSuite))
}
