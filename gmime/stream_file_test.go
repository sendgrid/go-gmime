package gmime_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/sendgrid/go-gmime/gmime"
	"github.com/stretchr/testify/assert"
)

//type FileStreamTestSuite AbstractStreamTestSuite

/*
func (s *FileStreamTestSuite) SetupTest() {
	s.buffer = "hello, world!"
    tempFilepath := pathJoin(os.TempDir(), "hello-gmime.txt")
    ioutil.WriteFile([]byte(s.buffer))
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
*/

func TestFileStreamWithPathTestSuite(t *testing.T) {
	s := "hello, world!"
	tempFilepath := path.Join(os.TempDir(), "hello-gmime.txt")

	{
		err := ioutil.WriteFile(tempFilepath, []byte(s), 0644)
		assert.NoError(t, err)
		defer os.Remove(tempFilepath)

		rfs := gmime.NewFileStreamForPath(tempFilepath, "r")
		l := rfs.Length()
		l2, r := rfs.Read(l)
		assert.Equal(t, int64(len(s)), l)
		assert.Equal(t, l, l2)
		assert.Equal(t, s, string(r))
	}

	{
		wfs := gmime.NewFileStreamForPath(tempFilepath, "w")
		defer os.Remove(tempFilepath)
		b := []byte(s)
		l := wfs.Write(b)
		assert.Equal(t, int64(len(b)), l)
		wfs.Close()

		r, err := ioutil.ReadFile(tempFilepath)
		assert.NoError(t, err)
		assert.Equal(t, s, string(r))
	}
	{
		wfs := gmime.NewFileStreamForPath(tempFilepath, "w")
		defer os.Remove(tempFilepath)
		l := wfs.WriteString(s)
		assert.Equal(t, int64(len(s)), l)
		wfs.Close()

		r, err := ioutil.ReadFile(tempFilepath)
		assert.NoError(t, err)
		assert.Equal(t, s, string(r))
	}
}
