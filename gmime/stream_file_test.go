package gmime

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FileStreamTestSuite AbstractStreamTestSuite

func (s *FileStreamTestSuite) SetupTest() {
	s.buffer = "hello, world!"
	tempFilepath := os.TempDir() + string(os.PathSeparator) + "hello.txt"
	tempFile, _ := os.Create(tempFilepath)
	assert.NotNil(s.T(), tempFile)
	tempFile.WriteString(s.buffer)
	tempFile.Sync()
	buffer := make([]byte, len(s.buffer))
	tempFile.Seek(0, os.SEEK_SET)
	tempFile.Read(buffer)
	assert.Equal(s.T(), s.buffer, string(buffer))
	tempFile.Seek(0, os.SEEK_SET)
	tempFile.Close()
}

func (s *FileStreamTestSuite) TearDownTest() {

}

func TestFileStreamTestSuite(t *testing.T) {
	suite.Run(t, new(FileStreamTestSuite))
}
