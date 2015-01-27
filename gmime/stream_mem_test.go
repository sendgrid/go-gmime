package gmime

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MemStreamTestSuite struct {
	suite.Suite
}

func (s *MemStreamTestSuite) TestNewMemStreamWithBuffer() {
	data, _ := ioutil.ReadFile("fixtures/text-only.eml")

	loop := 1
	for i := 0; i < loop; i++ {
		stream := NewMemStreamWithBuffer(string(data))
		assert.Equal(s.T(), stream.Length(), len(data))
	}
}

func (s *MemStreamTestSuite) TestBytes() {
	buffer := "hello, world!"
	stream := NewMemStreamWithBuffer(buffer)
	assert.Equal(s.T(), buffer, string(stream.Bytes()))
	stream.Close()
}

func TestMemStreamTestSuite(t *testing.T) {
	suite.Run(t, new(MemStreamTestSuite))
}
