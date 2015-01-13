package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MessagePartialTestSuite struct {
	suite.Suite
}

func (s *MessagePartialTestSuite) SetupTest() {
}

func (s *MessagePartialTestSuite) TearDownTest() {
}

func (s *MessagePartialTestSuite) TestNewMessagePart() {
	mp := NewMessagePartial("hola", 1, 1)
	assert.NotNil(s.T(), mp)
}

func (s *MessagePartialTestSuite) TestId() {
	mp := NewMessagePartial("hola", 1, 1)
	assert.Equal(s.T(), mp.Id(), "hola")
}

func (s *MessagePartialTestSuite) TestNumber() {
	mp := NewMessagePartial("hola", 1, 1)
	assert.Equal(s.T(), mp.Number(), 1)
}

func (s *MessagePartialTestSuite) TestTotal() {
	mp := NewMessagePartial("hola", 1, 1)
	assert.Equal(s.T(), mp.Total(), 1)
}

func (s *MessagePartialTestSuite) TestReconstructMessage() {
	// TODO: implement
}

func (s *MessagePartialTestSuite) TestSplitMessage() {
	// TODO: implement
}

// run test
func TestMessagePartialTestSuite(t *testing.T) {
	suite.Run(t, new(MessagePartialTestSuite))
}
