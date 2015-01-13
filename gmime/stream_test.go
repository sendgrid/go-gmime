package gmime

import (
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AbstractStreamTestSuite struct {
	suite.Suite
	buffer string
	stream Stream
}

func (s *AbstractStreamTestSuite) TestGenericStream() {
	buffer := s.buffer
	stream := s.stream

	assert.Equal(s.T(), stream.Length(), len(buffer))
	length, word := stream.Read((uint64)(5))
	assert.Equal(s.T(), 5, length)
	assert.Equal(s.T(), "hello", (string)(word))
	assert.Equal(s.T(), 5, stream.Tell())
	stream.Seek(5+2, os.SEEK_SET)
	assert.Equal(s.T(), 5+2, stream.Tell())
	assert.False(s.T(), stream.Eos())
	capitalW := [1]byte{(byte)('W')}
	stream.Write(([]byte)(capitalW[0:1]), (uint64)(1))
	stream.Reset()
	length, word = stream.Read((uint64)(len(buffer) + 1))
	assert.Equal(s.T(), len(buffer), length)
	assert.Equal(s.T(), "hello, World!", (string)(word))
	stream.Seek(0, os.SEEK_END)
	// XXX The File Stream doesn't know it's at EOF until another read is tried
	stream.Read((uint64)(1))
	assert.True(s.T(), stream.Eos())
	stream.Reset()
	subStream := stream.SubStream(0, 5)
	subString := "y"
	res := subStream.WriteString(subString)
	assert.Equal(s.T(), 1, res)
	subStream.Seek(-1, os.SEEK_CUR)
	length, word = subStream.Read((uint64)(4))
	assert.Equal(s.T(), "yell", string(word))
	outStream := NewMemStream()
	stream.WriteToStream(outStream)
	assert.Equal(s.T(), "yello, World!", string(outStream.Bytes()))
	outStream.Close()
}

// XXX Don't instantiate and run this test suite - it's an abstract base
