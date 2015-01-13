package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func (s *ParserTestSuite) TestParser() {
	message := NewMessage()
	textPart := NewPartWithType("text", "plain")
	htmlPart := NewPartWithType("text", "html")
	multipart := NewMultipart()
	multipart.AddPart(textPart)
	multipart.AddPart(htmlPart)
	message.SetMimePart(multipart)
	stream := NewMemStreamWithBuffer(message.ToString())
	defer stream.Close()

	parser := NewParserWithStream(stream)
	reconstructedMessage := parser.ConstructMessage()

	reconstructedMultipart, ok := reconstructedMessage.MimePart().(Multipart)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), "text/plain", reconstructedMultipart.GetPart(0).ContentType().ToString())
	assert.Equal(s.T(), "text/html", reconstructedMultipart.GetPart(1).ContentType().ToString())
}

// run test
func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}
