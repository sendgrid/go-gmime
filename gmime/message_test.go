package gmime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite
	message Message
}

func (s *MessageTestSuite) TestSender() {
	message := NewMessage()
	senderName := "hola"
	message.SetSender(senderName)
	sender := message.Sender()
	assert.Equal(s.T(), senderName, sender)

	// add 2nd sender, should still be 1st sender
	secondSender := "hola 2"
	message.SetSender(secondSender)
	sender = message.Sender()
	assert.Equal(s.T(), senderName, sender)
}

func (s *MessageTestSuite) TestReplyTo() {
	message := NewMessage()
	replyName := "hola"
	message.SetReplyTo(replyName)
	reply := message.ReplyTo()
	assert.Equal(s.T(), replyName, reply)

	// add 2nd replyTo, should not take it
	secondReply := ""
	message.SetReplyTo(secondReply)
	reply = message.ReplyTo()
	assert.Equal(s.T(), reply, secondReply)
}

func (s *MessageTestSuite) TestSubject() {
	subjectName := "hola"
	message := NewMessage()

	message.SetSubject(subjectName)
	subject := message.Subject()
	assert.Equal(s.T(), subjectName, subject)

	// add 2nd subject, should take it
	secondSubject := "hola 2"
	message.SetSubject(secondSubject)
	subject = message.Subject()
	assert.Equal(s.T(), secondSubject, subject)
}

func (s *MessageTestSuite) TestMessageId() {
	messageIdName := "hola"
	message := NewMessage()

	message.SetMessageId(messageIdName)
	messageId := message.MessageId()
	assert.Equal(s.T(), messageId, messageIdName)

	// add 2nd message id, should replace 1st id
	secondMessageIdName := "hola 2"
	message.SetMessageId(secondMessageIdName)
	messageId = message.MessageId()
	assert.Equal(s.T(), messageId, secondMessageIdName)
}

func (s *MessageTestSuite) TestTo() {
	message := NewMessage()
	firstName := "hola"
	firstEmail := "kickass@example.com"
	firstExpecting := fmt.Sprintf("%s <%s>", firstName, firstEmail)
	message.AddTo(firstName, firstEmail)
	info := message.To()
	assert.Equal(s.T(), info.ToString(true), firstExpecting)

	// add 2nd message id, should replace 1st id
	secondName := "hola2"
	secondEmail := "kickass2@example.com"
	secondExpecting := fmt.Sprintf("%s, %s <%s>", firstExpecting, secondName, secondEmail)
	message.AddTo(secondName, secondEmail)
	info = message.To()
	assert.Equal(s.T(), info.ToString(true), secondExpecting)
}

func (s *MessageTestSuite) TestCc() {
	message := NewMessage()
	firstName := "hola"
	firstEmail := "kickass@example.com"
	firstExpecting := fmt.Sprintf("%s <%s>", firstName, firstEmail)
	message.AddCc(firstName, firstEmail)
	info := message.Cc()
	assert.Equal(s.T(), info.ToString(true), firstExpecting)

	// add 2nd message id, should replace 1st id
	secondName := "hola2"
	secondEmail := "kickass2@example.com"
	secondExpecting := fmt.Sprintf("%s, %s <%s>", firstExpecting, secondName, secondEmail)
	message.AddCc(secondName, secondEmail)
	info = message.Cc()
	assert.Equal(s.T(), info.ToString(true), secondExpecting)
}

func (s *MessageTestSuite) TestBcc() {
	message := NewMessage()
	firstName := "hola"
	firstEmail := "kickass@example.com"
	firstExpecting := fmt.Sprintf("%s <%s>", firstName, firstEmail)
	message.AddBcc(firstName, firstEmail)
	info := message.Bcc()
	assert.Equal(s.T(), info.ToString(true), firstExpecting)

	// add 2nd message id, should replace 1st id
	secondName := "hola2"
	secondEmail := "kickass2@example.com"
	secondExpecting := fmt.Sprintf("%s, %s <%s>", firstExpecting, secondName, secondEmail)
	message.AddBcc(secondName, secondEmail)
	info = message.Bcc()
	assert.Equal(s.T(), info.ToString(true), secondExpecting)
}

func (s *MessageTestSuite) TestAllRecipients() {
	message := NewMessage()
	toName := "TO"
	toEmail := "to@example.com"
	toExpecting := fmt.Sprintf("%s <%s>", toName, toEmail)
	message.AddTo(toName, toEmail)

	ccName := "CC"
	ccEmail := "cc@example.com"
	ccExpecting := fmt.Sprintf("%s <%s>", ccName, ccEmail)
	message.AddCc(ccName, ccEmail)

	bccName := "BCC"
	bccEmail := "bcc@example.com"
	bccExpecting := fmt.Sprintf("%s <%s>", bccName, bccEmail)
	message.AddBcc(bccName, bccEmail)

	assert.Equal(s.T(), message.AllRecipients().GetLength(), 3)
	recipientsExpecting := fmt.Sprintf("%s, %s, %s", toExpecting, ccExpecting, bccExpecting)
	recipientsActual := message.AllRecipients().ToString(true)
	assert.Equal(s.T(), recipientsActual, recipientsExpecting)
}

func (s *MessageTestSuite) TestMimePart() {
	message := NewMessage()
	mimePart := NewPart()
	message.SetMimePart(mimePart)
	mime := message.MimePart()
	assert.NotNil(s.T(), mime)
	_, ok := mime.(Part)
	assert.True(s.T(), ok)
}

func (s *MessageTestSuite) TestMultiPart() {
	message := NewMessage()
	mimePart := NewMultipart()
	message.SetMimePart(mimePart)
	mime := message.MimePart()
	assert.NotNil(s.T(), mime)
	_, ok := mime.(Multipart)
	assert.True(s.T(), ok)
}

func (s *MessageTestSuite) TestToString() {
	message := NewMessage()
	senderName := "from"
	message.SetSender(senderName)
	replyTo := "reply"
	message.SetReplyTo(replyTo)
	subject := "subject"
	message.SetSubject(subject)
	messageId := "11111"
	message.SetMessageId(messageId)
	recipientName1 := "name 1"
	recipientEmail1 := "email1"
	recipientName2 := "name 2"
	recipientEmail2 := "email2"
	message.AddTo(recipientName1, recipientEmail1)
	message.AddTo(recipientName2, recipientEmail2)
	recipientName3 := "name 3"
	recipientEmail3 := "email3"
	recipientName4 := "name 4"
	recipientEmail4 := "email4"
	message.AddCc(recipientName3, recipientEmail3)
	message.AddBcc(recipientName4, recipientEmail4)

	text := "This is a text part"
	textStream := NewMemStreamWithBuffer(text)
	textEncoding := NewContentEncodingFromString("8bit")
	textWrapper := NewDataWrapperWithStream(textStream, textEncoding)
	textPart := NewPartWithType("text", "plain")
	textPart.SetContentObject(textWrapper)

	message.SetMimePart(textPart)
	mimeMessageActual := message.ToString()

	mimeMessageExpected := fmt.Sprintf(
		`From: %s
Reply-To: %s
Subject: %s
Message-Id: <%s>
To: %s <%s>, %s <%s>
Cc: %s <%s>
Bcc: %s <%s>
MIME-Version: 1.0
Content-Type: %s

%s`,
		senderName, replyTo, subject, messageId,
		recipientName1, recipientEmail1, recipientName2, recipientEmail2,
		recipientName3, recipientEmail3, recipientName4, recipientEmail4,
		textPart.ContentType().ToString(), text)

	assert.Equal(s.T(), mimeMessageActual, mimeMessageExpected)

	html := "<html><body>This is an HTML part</body></hmtl>"
	htmlStream := NewMemStreamWithBuffer(html)
	htmlEncoding := NewContentEncodingFromString("8bit")
	htmlWrapper := NewDataWrapperWithStream(htmlStream, htmlEncoding)
	htmlPart := NewPartWithType("text", "html")
	htmlPart.SetContentObject(htmlWrapper)

	multipart := NewMultipart()
	multipart.AddPart(textPart)
	multipart.AddPart(htmlPart)

	message.SetMimePart(multipart)
	multipartMessageActual := message.ToString()
	boundary := message.MimePart().ContentType().Parameter("boundary")

	multipartMessageExpected := fmt.Sprintf(
		`From: %s
Reply-To: %s
Subject: %s
Message-Id: <%s>
To: %s <%s>, %s <%s>
Cc: %s <%s>
Bcc: %s <%s>
MIME-Version: 1.0
Content-Type: %s; boundary="%s"

--%s
Content-Type: %s

%s
--%s
Content-Type: %s

%s
--%s--
`,
		senderName, replyTo, subject, messageId,
		recipientName1, recipientEmail1, recipientName2, recipientEmail2,
		recipientName3, recipientEmail3, recipientName4, recipientEmail4,
		multipart.ContentType().ToString(), boundary,
		boundary, textPart.ContentType().ToString(), text,
		boundary, htmlPart.ContentType().ToString(), html,
		boundary)

	assert.Equal(s.T(), multipartMessageActual, multipartMessageExpected)
}

func (s *MessageTestSuite) TestContentType() {
	message := NewMessage()
	typename := "text/plain"
	contentType := NewContentTypeFromString(typename)
	message.SetContentType(contentType)
	assert.Equal(s.T(), message.ContentType().ToString(), typename)
}

// run test
func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}
