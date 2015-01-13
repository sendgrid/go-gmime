package gmime

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	assert "github.com/stretchr/testify/assert"
	suite "github.com/stretchr/testify/suite"
	"os"
	"runtime/debug"
	"testing"
)

type ParseMessageTestSuite struct {
	suite.Suite
}

//Example on how to create a NewParse
func ExampleNewParse() {
	fileHandler, _ := os.Open("fixtures/text_attachment.eml")
	defer fileHandler.Close()
	reader := bufio.NewReader(fileHandler)
	parse := NewParse(reader)

	fmt.Println(parse.From())
	fmt.Println(parse.To())
	fmt.Println(parse.Subject())

	// Output:
	// Dave McGuire <foobar@foobar.com>
	// Team R&D <foobar@foobar.com>
	// Distilled Content-Type headers (outbound & failed)
}

func (s *ParseMessageTestSuite) TestParseTextOnly() {
	loop := 1

	for i := 0; i < loop; i++ {
		fileHandler, err := os.Open("fixtures/text-only.eml")
		assert.Nil(s.T(), err)

		reader := bufio.NewReader(fileHandler)
		parse := NewParse(reader)

		assert.Equal(s.T(), parse.To(), "Kien Pham <kien@sendgrid.com>")
		assert.Equal(s.T(), parse.Subject(), "text only")
		assert.Equal(s.T(), parse.From(), "Kien Pham <kien.pham@sendgrid.com>")
		assert.Equal(s.T(), parse.Bcc(), "")
		assert.Equal(s.T(), parse.Cc(), "")
		assert.Equal(s.T(), parse.Recipients(), "Kien Pham <kien@sendgrid.com>")
		assert.Equal(s.T(), parse.MessageId(), "CAGPJ=uZ8BfOdJr9E-J3o=5uC=4j0YECrm6Aa58d8vovNNrMS4Q@mail.gmail.com")
		assert.Equal(s.T(), parse.Text(), "this is text only email")
		assert.Equal(s.T(), parse.Html(), "")
		attachments := parse.Attachment()
		assert.Equal(s.T(), len(attachments), 0)
		fileHandler.Close()
	}
}

func (s *ParseMessageTestSuite) TestParseMultipart() {
	loop := 1

	for i := 0; i < loop; i++ {
		fileHandler, err := os.Open("fixtures/multipart-alternative.eml")
		assert.Nil(s.T(), err)

		reader := bufio.NewReader(fileHandler)
		parse := NewParse(reader)

		assert.Equal(s.T(), parse.To(), "joe@sixpack.org")
		assert.Equal(s.T(), parse.Subject(), "Test")
		assert.Equal(s.T(), parse.From(), "Jihad <jihad@vt.edu>")
		assert.Equal(s.T(), parse.Bcc(), "")
		assert.Equal(s.T(), parse.Cc(), "")
		assert.Equal(s.T(), parse.Recipients(), "joe@sixpack.org")
		assert.Equal(s.T(), parse.MessageId(), "000901bf1857$25c23850$66d9c026@jackhandy")

		hasher := md5.New()
		hasher.Write([]byte(parse.Text()))
		hashString := hex.EncodeToString(hasher.Sum(nil))
		assert.Equal(s.T(), hashString, "28e0cd851e8c6a443813c6178dc61213")
		hasher.Reset()

		hasher.Write([]byte(parse.Html()))
		hashString = hex.EncodeToString(hasher.Sum(nil))
		assert.Equal(s.T(), hashString, "e9c42f1ed2abfc23603896e2e8c31568")
		hasher.Reset()

		attachments := parse.Attachment()
		for key, value := range attachments {
			assert.Equal(s.T(), key, "leonc.jpg")

			hasher.Write(value)
			hashString = hex.EncodeToString(hasher.Sum(nil))
			assert.Equal(s.T(), hashString, "44dcb9ed2fb0e046eab2913a9ecccace")
			hasher.Reset()
		}
		fileHandler.Close()
	}
}

func (s *ParseMessageTestSuite) TestParseLargeAttachments() {
	loop := 1
	for i := 0; i < loop; i++ {
		fileHandler, err := os.Open("fixtures/large_attachments.eml")
		assert.Nil(s.T(), err)

		reader := bufio.NewReader(fileHandler)
		parse := NewParse(reader)

		assert.Equal(s.T(), parse.To(), "Kien Pham <kien@sendgrid.com>")
		assert.Equal(s.T(), parse.Subject(), "break my parse test!!!")
		assert.Equal(s.T(), parse.From(), "Kien Pham <kien.pham@sendgrid.com>")
		assert.Equal(s.T(), parse.Bcc(), "")
		assert.Equal(s.T(), parse.Cc(), "")
		assert.Equal(s.T(), parse.Recipients(), "Kien Pham <kien@sendgrid.com>")
		assert.Equal(s.T(), parse.MessageId(), "CAGPJ=ubOouY1J7wqDEMKdTVzs3xYg7xa=UcupvZzWCGk0nMP1Q@mail.gmail.com")

		hasher := md5.New()
		hasher.Write([]byte(parse.Text()))
		hashString := hex.EncodeToString(hasher.Sum(nil))
		assert.Equal(s.T(), hashString, "3a84d07f16dd41d574054e3762b5fb0e")
		hasher.Reset()

		hasher.Write([]byte(parse.Html()))
		hashString = hex.EncodeToString(hasher.Sum(nil))
		assert.Equal(s.T(), hashString, "d864471b159bafe6271c40c7dddb238b")
		hasher.Reset()

		attachmentName := make(map[int]string)
		attachmentName[0] = "DSC_0103 (2).JPG"
		attachmentName[1] = "DSC_0102 (2).JPG"

		attachmentHash := make(map[string]string)
		attachmentHash[attachmentName[0]] = "fce92d641bbf8c8ed7152646a1d6b3e7"
		attachmentHash[attachmentName[1]] = "981ed617e62597012d54e81a21420e1d"

		attachments := parse.Attachment()
		_ = attachments
		counter := 0
		for key, value := range attachments {
			//TODO: kills valgrind
			//assert.Equal(s.T(), key, attachmentName[counter])
			assert.NotNil(s.T(), value)

			hasher.Write(value)
			hashString = hex.EncodeToString(hasher.Sum(nil))
			assert.Equal(s.T(), hashString, attachmentHash[key])
			hasher.Reset()
			counter++
		}
		fileHandler.Close()
	}
}

func (s *ParseMessageTestSuite) TestParseTextAttachment() {
	loop := 1
	for i := 0; i < loop; i++ {
		fileHandler, err := os.Open("fixtures/text_attachment.eml")
		assert.Nil(s.T(), err)

		reader := bufio.NewReader(fileHandler)
		parse := NewParse(reader)
		assert.Equal(s.T(), parse.To(), "Team R&D <foobar@foobar.com>")
		hasher := md5.New()
		hasher.Write([]byte(parse.Text()))
		hashString := hex.EncodeToString(hasher.Sum(nil))
		assert.Equal(s.T(), hashString, "b5b94cd495174ab6e4443fa81847b6ce")
		hasher.Reset()

		hasher.Write([]byte(parse.Html()))
		hashString = hex.EncodeToString(hasher.Sum(nil))
		assert.Equal(s.T(), hashString, "d41d8cd98f00b204e9800998ecf8427e")
		hasher.Reset()

		attachments := parse.Attachment()
		for key, value := range attachments {
			assert.Equal(s.T(), key, "content-type-sorted-distilled.txt")

			hasher.Write(value)
			hashString = hex.EncodeToString(hasher.Sum(nil))
			assert.Equal(s.T(), hashString, "bcd2881b9b06f24960d4998567be37bd")
			hasher.Reset()
		}
		fileHandler.Close()
	}
}

func (s *ParseMessageTestSuite) TestInlineImageAttachment() {
	loop := 1
	for i := 0; i < loop; i++ {
		fileHandler, err := os.Open("fixtures/inline-attachment.eml")
		assert.Nil(s.T(), err)

		reader := bufio.NewReader(fileHandler)
		parse := NewParse(reader)
		hasher := md5.New()

		attachments := parse.Attachment()

		for key, value := range attachments {
			assert.Equal(s.T(), key, "kien.jpg")

			hasher.Write(value)
			hashString := hex.EncodeToString(hasher.Sum(nil))
			assert.Equal(s.T(), hashString, "322fb6caca84ddddea8d62c0e3b85d8e")
			hasher.Reset()
		}
		fileHandler.Close()
	}
	debug.FreeOSMemory()
}

// run test
func TestParseMessageTestSuite(s *testing.T) {
	suite.Run(s, new(ParseMessageTestSuite))
	debug.FreeOSMemory()
}
