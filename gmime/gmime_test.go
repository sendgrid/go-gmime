package gmime_test

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"strings"

	"github.com/sendgrid/go-gmime/gmime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingAMessageWithMultipartAlternativeMIMEType(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/multipart-alternative.eml")
	assert.NoError(t, err)

	stream := gmime.NewMemStreamWithBuffer(string(data))
	parser := gmime.NewParserWithStream(stream)
	message := parser.ConstructMessage()

	// "should have a hash that matches the hash of the file on disk"
	{
		h := sha1.New()
		h.Write([]byte(message.ToString()))
		assert.Equal(t, hex.EncodeToString(h.Sum(nil)), "f47e39511cd3be8630d67098266f6a3727919a96")
	}
	// "should show the message mime part is a multipart/mixed container"
	{
		container, ok := message.MimePart().(gmime.Multipart)
		assert.True(t, ok)
		assert.NotNil(t, container)
		contentType := strings.ToLower(container.ContentType().ToString())
		assert.Equal(t, contentType, "multipart/mixed")
	}

	// "When the iterator on the message is run"
	{
		parts := make([]gmime.Part, 0)
		multiparts := make([]gmime.Multipart, 0)

		for iter := gmime.NewPartIter(message); iter.HasNext(); iter.Next() {
			object := iter.Current()

			if part, ok := object.(gmime.Multipart); ok {
				multiparts = append(multiparts, part)
			} else if part, ok := object.(gmime.Part); ok {
				parts = append(parts, part)
			}
		}

		assert.Equal(t, len(multiparts), 2,
			"should return a non-empty multipart container list")

		assert.Equal(t, len(parts), 3,
			"should return as many items as the number of parts on disk")

		//		"When the parts of the message are examined"
		// "should show the message mime part is a multipart/mixed container"
		{
			container := multiparts[0]
			assert.NotNil(t, container)
			contentType := strings.ToLower(container.ContentType().ToString())
			assert.Equal(t, contentType, "multipart/mixed",
				"should show the message mime part is a multipart/mixed container")
		}

		{
			contentType := strings.ToLower(multiparts[1].ContentType().ToString())
			assert.Equal(t, contentType, "multipart/alternative",
				"should have a multipart container that matches the one on disk")
		}

		{
			// "should have a text part that matches the one on disk"
			part := parts[0]
			contentType := strings.ToLower(part.ContentType().ToString())
			assert.Equal(t, contentType, "text/plain")
			size, hash := sizeAndHashOf(part)
			assert.Equal(t, size, 499)
			assert.Equal(t, hash, "28e0cd851e8c6a443813c6178dc61213")
		}

		{
			// "should have an html part that matches the one on disk"
			part := parts[1]
			contentType := strings.ToLower(part.ContentType().ToString())
			assert.Equal(t, contentType, "text/html")
			size, hash := sizeAndHashOf(part)
			assert.Equal(t, size, 2173)
			assert.Equal(t, hash, "e9c42f1ed2abfc23603896e2e8c31568")
		}

		{
			// "should have an image part that matches the one on disk"
			part := parts[2]
			contentType := strings.ToLower(part.ContentType().ToString())
			assert.Equal(t, contentType, "image/jpeg")
			size, hash := sizeAndHashOf(part)
			assert.Equal(t, size, 17527)
			assert.Equal(t, hash, "44dcb9ed2fb0e046eab2913a9ecccace")
		}
	}
}

// "Parsing a message with an RFC822 message MIME type as a part",
func TestParsingAMessageWithAnRFC822(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/message-as-part.eml")
	assert.NoError(t, err)

	stream := gmime.NewMemStreamWithBuffer(string(data))
	parser := gmime.NewParserWithStream(stream)
	message := parser.ConstructMessage()

	// "When the message has been constructed"
	//	  "should have a hash that matches the hash of the file on disk"
	{
		h := sha1.New()
		h.Write([]byte(message.ToString()))
		assert.Equal(t, hex.EncodeToString(h.Sum(nil)), "2ae3c8c00e9389a287d7ace9e7584d96fcffba01")
	}

	// "When the iterator on the message is run"
	{
		var multipart gmime.Multipart
		parts := make([]gmime.Part, 0)
		var rfc822 gmime.MessagePart

		for iter := gmime.NewPartIter(message); iter.HasNext(); iter.Next() {
			object := iter.Current()

			if part, ok := object.(gmime.Multipart); ok {
				multipart = part
			} else if part, ok := object.(gmime.Part); ok {
				parts = append(parts, part)
			} else if part, ok := object.(gmime.MessagePart); ok {
				rfc822 = part
			}
		}

		// "should show the message mime part is a multipart/mixed container"
		{
			container := multipart
			assert.NotNil(t, container)
			contentType := strings.ToLower(container.ContentType().ToString())
			assert.Equal(t, contentType, "multipart/mixed")
		}

		// "should return a non-nil RFC822 container as one of the parts"
		assert.NotNil(t, rfc822)

		// "should return as many items as the number of parts on disk"
		assert.Equal(t, len(parts), 2)

		// "When the parts of the message are examined"
		{
			// "should have a text part that matches the one on disk"
			{
				part := parts[0]
				contentType := strings.ToLower(part.ContentType().ToString())
				assert.Equal(t, contentType, "text/plain")
				size, hash := sizeAndHashOf(part)
				assert.Equal(t, size, 262)
				assert.Equal(t, hash, "b5aec669cdae7ef2b62ef0b3070c8fd7")
			}

			// "should have an RFC822 container that matches the one on disk"
			{
				contentType := strings.ToLower(rfc822.ContentType().ToString())
				assert.Equal(t, contentType, "message/rfc822")
			}

			// "should have an *embedded* text part that matches the one on disk"
			{
				part := parts[1]
				contentType := strings.ToLower(part.ContentType().ToString())
				assert.Equal(t, contentType, "text/plain")
				size, hash := sizeAndHashOf(part)
				assert.Equal(t, size, 1619)
				assert.Equal(t, hash, "9fb795da046d6a5338a5ade3ebe3c1b1")
			}
		}
	}
}

// "Parsing a Delivery Status Notification (DSN) for a bounce"
func ParsingADeliveryStatusNotification(t *testing.T) {
	data, err := ioutil.ReadFile("fixtures/DSN-bounce.eml")
	assert.NoError(t, err)

	stream := gmime.NewMemStreamWithBuffer(string(data))
	parser := gmime.NewParserWithStream(stream)
	message := parser.ConstructMessage()

	// "When the message has been constructed"
	{
		// "should have a hash that matches the hash of the file on disk"
		{
			h := sha1.New()
			h.Write([]byte(message.ToString()))
			assert.Equal(t, hex.EncodeToString(h.Sum(nil)), "fca9901a01ab33448c0b8c5f64f9f19404560234")
		}
	}

	// "When the iterator on the message is run"
	{
		var report gmime.Multipart
		var preamble gmime.Part
		var status gmime.Part
		var rfc822 gmime.MessagePart
		var ok bool

		iter := gmime.NewPartIter(message)

		// "should show the message mime part is a multipart/report container"
		{
			report, ok = iter.Current().(gmime.Multipart)
			assert.True(t, ok)
			container := report

			assert.NotNil(t, container)
			contentType := strings.ToLower(container.ContentType().ToString())
			assert.Equal(t, contentType, "multipart/report")
		}

		iter.Next()

		// "should have a human-readable preamble for the bounced message"
		{
			preamble, ok = iter.Current().(gmime.Part)
			assert.True(t, ok)
			assert.NotNil(t, preamble)
			contentType := strings.ToLower(preamble.ContentType().ToString())
			assert.Equal(t, contentType, "text/plain")
			assert.Equal(t, preamble.Header("Content-Description"), "Notification")
			explanation := preamble.ToString()
			assert.Contains(t, explanation, "not able to be")
			assert.Contains(t, explanation, "delivered to one of its intended recipients")
			assert.Contains(t, explanation, "550 5.1.1 sid=i01K1n00l0kn1Em01 Address rejected tobigeri-555@mail.goo.ne.jp. [code=28]")
		}

		iter.Next()

		// "should have a machine-readable message/delivery-status for the bounced message"
		{
			status, ok = iter.Current().(gmime.Part)
			assert.True(t, ok)
			assert.NotNil(t, status)
			contentType := strings.ToLower(status.ContentType().ToString())
			assert.Equal(t, contentType, "message/delivery-status")
			assert.Equal(t, status.Header("Action"), "failed")
			assert.Equal(t, status.Header("Status"), "5.1.1")
			assert.Contains(t, status.ToString(), "Arrival-Date: 2014-03-26 00-01-19")
			assert.Contains(t, status.ToString(), "Diagnostic-Code: 550 5.1.1 sid=i01K1n00l0kn1Em01 Address rejected tobigeri-555@mail.goo.ne.jp. [code=28]")
		}

		iter.Next()

		// "should return a non-nil RFC822 container as one of the parts"
		{
			rfc822, ok = iter.Current().(gmime.MessagePart)
			assert.True(t, ok)
			assert.NotNil(t, rfc822)
			contentType := strings.ToLower(rfc822.ContentType().ToString())
			assert.Equal(t, contentType, "message/rfc822")
			assert.Equal(t, rfc822.Header("Content-Description"), "Undelivered Message")
		}

		// "When examining the container part of the original message"
		{
			original := rfc822.Message()
			assert.NotNil(t, original)
			part := original.MimePart()
			contentType := strings.ToLower(part.ContentType().ToString())
			assert.Equal(t, contentType, "multipart/alternative")

			// "should have enough information to be able to record the bounce"
			{
				assert.Equal(t, original.Header("To"), "tobigeri-555@mail.goo.ne.jp")
				// These headers encode things like user id, etc.
				assert.NotNil(t, original.Header("X-SG-EID"))
				assert.NotNil(t, original.Header("X-SG-ID"))
			}
		}
	}
}

func sizeAndHashOf(p gmime.Part) (int, string) {
	dataWrapper := p.ContentObject()
	memStream := gmime.NewMemStream()
	c := dataWrapper.WriteToStream(memStream)
	memStream.Flush()
	data := memStream.Bytes()
	hasher := md5.New()
	hasher.Write(data)
	hashString := hex.EncodeToString(hasher.Sum(nil))
	return int(c), hashString
}
