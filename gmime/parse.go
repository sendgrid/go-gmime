package gmime

import (
	"io"
	"io/ioutil"
	"strings"

	iconv "github.com/djimenez/iconv-go"
)

type Parse interface {
	To() string
	Subject() (string, bool)
	From() (string, bool)
	Bcc() string
	Cc() string
	Recipients() string
	MessageId() (string, bool)
	Text() (string, bool)
	Html() (string, bool)
	Attachment() map[string][]byte
	//	AttachmentReader() Reader
	Headers() string
	Header(string) (string, bool)
	Boundary() string
}

type aParse struct {
	message Message
}

func NewParse(reader io.Reader) Parse {
	// TODO: stream instead read all in at once
	data, _ := ioutil.ReadAll(reader)
	stream := NewMemStreamWithBuffer(string(data))
	parser := NewParserWithStream(stream)
	message := parser.ConstructMessage()
	return &aParse{
		message: message,
	}
}

func (p *aParse) To() string {
	return p.message.To().ToString(true)
}

func (p *aParse) Subject() (string, bool) {
	return p.message.Subject()
}

func (p *aParse) From() (string, bool) {
	return p.message.Sender()
}

func (p *aParse) Bcc() string {
	return p.message.Bcc().ToString(true)
}

func (p *aParse) Cc() string {
	return p.message.Cc().ToString(true)
}

func (p *aParse) Recipients() string {
	return p.message.AllRecipients().ToString(true)
}

func (p *aParse) MessageId() (string, bool) {
	return p.message.MessageId()
}

func (p *aParse) ContentType() string {
	return p.message.ContentType().ToString()
}

func (p *aParse) Text() (string, bool) {
	return p.parseBody("text/plain")
}

func (p *aParse) Html() (string, bool) {
	return p.parseBody("text/html")
}

func (p *aParse) Attachment() map[string][]byte {
	payload := make(map[string][]byte)

	container := p.message.MimePart()

	if part, ok := container.(Part); ok {
		if cd := part.ContentDisposition(); cd != nil {
			if cd.IsAttachment() {
				payload[part.Filename()] = p.rawPart(part)
			}
		}
	} else if _, ok := container.(Multipart); ok {
		for iter := NewPartIter(p.message); iter.HasNext(); iter.Next() {
			object := iter.Current()
			if part, ok = object.(Part); ok {
				if cd := part.ContentDisposition(); cd != nil {
					if cd.IsAttachment() {
						payload[part.Filename()] = p.rawPart(part)
					}
				}
			}
		}
	}
	return payload
}

//func (p *aParse) AttachmentReader() Reader {
//	// TODO: implement
//	return p.reader
//}

func (p *aParse) Headers() string {
	return p.message.Headers()
}

func (p *aParse) Header(name string) (string, bool) {
	return p.message.Header(name)
}

func (p *aParse) Boundary() string {
	payload := ""
	container := p.message.MimePart()
	if multipart, ok := container.(Multipart); ok {
		payload = multipart.Boundary()
	}

	return payload
}

// return the byte content of the Part. Ex: attachment
func (p *aParse) rawPart(part Part) []byte {
	writeStream := NewMemStream()
	if dataWrapper := part.ContentObject(); dataWrapper != nil {
		dataWrapper.WriteToStream(writeStream)
	}
	writeStream.Flush()

	return writeStream.Bytes()
}

// parse a Part
func (p *aParse) parsePart(part Part, contentType string) (string, bool) {
	if part.ContentType().ToString() == contentType && part.Filename() == "" {
		payload := string(p.rawPart(part))

		// convert charset
		targetCharset := "utf-8"
		sourceCharset := strings.ToLower(part.ContentType().Parameter("charset"))
		if sourceCharset != targetCharset {
			payload, _ = iconv.ConvertString(payload, sourceCharset, targetCharset)
		}
        return payload, true
	}

	return "", false
}

// parse the message body, might contains many Parts
func (p *aParse) parseBody(contentType string) (string, bool) {
	container := p.message.MimePart()
    payload := ""

	if part, ok := container.(Part); ok {
		return p.parsePart(part, contentType)
	} else if _, ok := container.(Multipart); ok {
		for iter := NewPartIter(p.message); iter.HasNext(); iter.Next() {
			if object := iter.Current(); object != nil {
				if part, ok = object.(Part); ok {
					// TODO: looks wrong
                    if payloadPart, ok := p.parsePart(part, contentType); ok {
                        payload += payloadPart
                    }
				}
			}
		}
	}

	return payload, len(payload) > 0
}
