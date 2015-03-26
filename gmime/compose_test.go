package gmime

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"runtime/debug"
)

type ComposerTestSuite struct {
	suite.Suite
	SenderEmail          string
	SenderEmailWithName  string
	ReceiverEmail        string
	ReceiverName         string
	ReceiverNameAndEmail string
	TextBody             string
	HTMLBody             string
	Subject              string
	ThisFileHandler      *os.File
}

// this run before each test
func (s *ComposerTestSuite) SetupTest() {
	s.SenderEmail = "senderEmail@example.com"
	s.SenderEmailWithName = "Some Sender <sender@example.com>"
	s.ReceiverEmail = "receiver@example.com"
	s.ReceiverName = "Some Receiver"
	s.ReceiverNameAndEmail = fmt.Sprintf("%s <%s>", s.ReceiverName, s.ReceiverEmail)
	s.TextBody = "This is a simple content of a text email"
	s.HTMLBody = "This is a very special <strong>HTML</strong> message <br /><br /> Thanks!"
	s.Subject = "This is message subject"

	// open this test file:
	s.ThisFileHandler, _ = os.Open("composer_test.go")
}

func (s *ComposerTestSuite) TearDownTest() {
	s.ThisFileHandler.Close()
}

//Simple message from sender@example.com to awesome@example.com
func ExampleNewCompose() {
	composer := NewCompose()
	composer.AddFrom("Good Sender <good_sender@example.com>")
	composer.AddTo("good_customer@example.com", "Good Customer")

	// read data from a file:
	fileHandler, _ := os.Open("composer_test.go")
	defer fileHandler.Close()
	reader := bufio.NewReader(fileHandler)

	composer.AddHTMLReader(reader)
}

func (s *ComposerTestSuite) TestNewTextComposer() {
	loop := 1
	for i := 0; i < loop; i++ {
		composer := NewCompose()
		composer.AddFrom(s.SenderEmail)
		composer.AddTo(s.ReceiverName, s.ReceiverEmail)
		reader := bufio.NewReader(s.ThisFileHandler)
		composer.AddTextReader(reader)

		from, ok := composer.From()
		assert.True(s.T(), ok)
		assert.Equal(s.T(), from, s.SenderEmail)
		assert.Equal(s.T(), composer.To(), s.ReceiverNameAndEmail)

		// TODO: test this
		//	assert.Equal(s.T(), composer.Text(), content of this file...)
	}
}

func (s *ComposerTestSuite) TestNewHTMLComposer() {
	loop := 1
	for i := 0; i < loop; i++ {
		composer := NewCompose()
		composer.AddFrom(s.SenderEmail)
		composer.AddTo(s.ReceiverName, s.ReceiverEmail)
		reader := bufio.NewReader(s.ThisFileHandler)
		composer.AddHTMLReader(reader)
		// TODO: test this
		//	assert.Equal(s.T(), composer.HTML(), content of this file...)
	}
}

func (s *ComposerTestSuite) TestNewMixedComposer() {
	loop := 1
	for i := 0; i < loop; i++ {
		composer := NewCompose()
		composer.AddFrom(s.SenderEmail)
		composer.AddTo(s.ReceiverName, s.ReceiverEmail)
		composer.AddCC(s.ReceiverName, s.ReceiverEmail)
		composer.AddBCC(s.ReceiverName, s.ReceiverEmail)
		composer.AddSubject(s.Subject)

		composer.AddText("Text part # 1", "")
		composer.AddHTML("<strong>HTML</strong> part #1", "8bit")
		composer.AddText("Text part # 2", "8bit")
		composer.AddHTML("<strong>HTML</strong> part #2", "8bit")

		expected := `From: senderEmail@example.com
To: Some Receiver <receiver@example.com>
Cc: Some Receiver <receiver@example.com>
Bcc: Some Receiver <receiver@example.com>
Subject: This is message subject
MIME-Version: 1.0
Content-Type: multipart/mixed

Content-Type: text/plain
Text part # 1

Content-Type: text/html
<strong>HTML</strong> part #1

Content-Type: text/plain
Text part # 2

Content-Type: text/html
<strong>HTML</strong> part #2

`
		assert.Equal(s.T(), composer.ToString(), expected)

		// TODO: test for AddTextReader, AddHTMLReader
		//		reader := bufio.NewReader(s.ThisFileHandler)
		//		composer.AddTextReader(reader)
		//		composer.AddHTMLReader(reader)
	}
}

// run test
func TestComposerTestSuite(t *testing.T) {
	suite.Run(t, new(ComposerTestSuite))
	debug.FreeOSMemory()
}
