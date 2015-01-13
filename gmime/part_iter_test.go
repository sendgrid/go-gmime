package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PartIterTestSuite struct {
	suite.Suite
}

func (s *PartIterTestSuite) TestNewPartIterOnMimeMessage() {
	message := NewMessage()
	htmlPart := NewPartWithType("text", "html")
	types := [1]string{"text/html"}
	message.SetMimePart(htmlPart)
	i := 0

	for partIter := NewPartIter(message); partIter.HasNext(); partIter.Next() {
		assert.Equal(s.T(), types[i], partIter.Current().ContentType().ToString())
		i++
	}

	assert.Equal(s.T(), len(types), i)
}

func (s *PartIterTestSuite) TestNewPartIterOnMultipartMessage() {
	message := NewMessage()
	textPart := NewPartWithType("text", "plain")
	htmlPart := NewPartWithType("text", "html")
	types := [3]string{"multipart/mixed", "text/plain", "text/html"}
	multipart := NewMultipart()
	multipart.AddPart(textPart)
	multipart.AddPart(htmlPart)
	message.SetMimePart(multipart)
	i := 0

	for partIter := NewPartIter(message); partIter.HasNext(); partIter.Next() {
		assert.Equal(s.T(), types[i], partIter.Current().ContentType().ToString())
		i++
	}

	assert.Equal(s.T(), len(types), i)
}

func TestPartIterTestSuite(t *testing.T) {
	suite.Run(t, new(PartIterTestSuite))
}
