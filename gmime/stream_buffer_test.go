package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"runtime/debug"
)

type StreamBufferTestSuite struct {
	suite.Suite
	stream Stream
}

// this run before each test
func (s *StreamBufferTestSuite) SetupTest() {
	s.stream = NewMemStream()
}

func (s *StreamBufferTestSuite) TestNewBufferStream() {
	cacheReadBufferStream := NewBufferedStream(s.stream, CACHE_READ)
	assert.Equal(s.T(), cacheReadBufferStream.Length(), 0)

	blockReadBufferStream := NewBufferedStream(s.stream, BLOCK_READ)
	assert.Equal(s.T(), blockReadBufferStream.Length(), 0)

	blockWriteBufferStream := NewBufferedStream(s.stream, BLOCK_WRITE)
	assert.Equal(s.T(), blockWriteBufferStream.Length(), 0)

	// TODO: stream some data
}

// TODO: implement test
func (s *StreamBufferTestSuite) TestBufferStreamGets() {

}

// TODO: implement test
func (s *StreamBufferTestSuite) TestBufferStreamReadLn() {

}

// run test
func TestStreamBufferTestSuite(t *testing.T) {
	suite.Run(t, new(StreamBufferTestSuite))
	debug.FreeOSMemory()
}
