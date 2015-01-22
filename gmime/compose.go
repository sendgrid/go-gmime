package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"
import (
	"io"
	"strings"
)

type Compose interface {
	AddTo(string, string)
	AddCC(string, string)
	AddBCC(string, string)
	AddSubject(string)
	AddFrom(string)
	AddText(string, string)
	AddHTML(string, string)
	AddTextReader(io.Reader)
	AddHTMLReader(io.Reader)

	To() string
	Subject() string
	Text() string
	HTML() string
	From() string
	Recipient() string
	ToString() string

	// TODO: implements
	//	AddAttachment(filePath string)
	//	AddAttachmentLink(filename string, url string)
	//	AddAttachmentStream(filename string, content []byte)
}

type aCompose struct {
	message    Message
	multipart  Multipart
	partsCount int
}

func NewCompose() Compose {
	return &aCompose{
		message:   NewMessage(),
		multipart: NewMultipartWithSubtype("mixed")}
}

func (p *aCompose) AddTo(name string, email string) {
	p.message.AddTo(name, email)
}

func (p *aCompose) AddCC(name string, email string) {
	p.message.AddCc(name, email)
}

func (p *aCompose) AddBCC(name string, email string) {
	p.message.AddBcc(name, email)
}

func (p *aCompose) AddSubject(subject string) {
	p.message.SetSubject(subject)
}

func (p *aCompose) AddText(text string, encoding string) {
	contentType := "text/plain"
	p.addPart(text, contentType, encoding)
}

func (p *aCompose) AddTextReader(reader io.Reader) {
	// TODO: implement reader

}

func (p *aCompose) AddHTML(html string, encoding string) {
	contentType := "text/html"
	p.addPart(html, contentType, encoding)
}

func (p *aCompose) AddHTMLReader(reader io.Reader) {
	// TODO: implement
}

func (p *aCompose) AddFrom(from string) {
	p.message.SetSender(from)
}

func (p *aCompose) Message() string {
	return p.message.ToString()
}

func (p *aCompose) To() string {
	return p.message.To().ToString(true)
}

func (p *aCompose) Subject() string {
	return p.message.Subject()
}

func (p *aCompose) Text() string {
	// TODO: implement?
	return ""
}

func (p *aCompose) HTML() string {
	// TODO: implement get HTML body part
	return p.Text()
}

func (p *aCompose) From() string {
	return p.message.Sender()
}

func (p *aCompose) Recipient() string {
	recipient := p.message.AllRecipients()

	return recipient.ToString(true)
}

func (p *aCompose) ToString() string {
	// write header
	payload := p.message.Headers() + "\n"

	part := p.message.MimePart()
	payload += p.part(part)

	return payload
}

func (p *aCompose) part(part Object) string {
	payload := ""

	switch part.(type) {
	case Multipart:
		m := part.(Multipart)
		n := m.Count()
		for i := 0; i < n; i++ {
			subpart := m.GetPart(i)
			payload += p.part(subpart)
		}

	case MessagePart:
		// TODO: implement
		println("MessagePart----------")

	case Part:
		m := part.(Part)
		dw := m.ContentObject()
		resultStream := dw.Stream()
		resultStream.Flush()
		_, data := resultStream.Read(resultStream.Length())

		if p.partsCount > 1 {
			payload = m.Headers()
		}

		payload += string(data) + "\n\n"
	}

	return payload
}

func (p *aCompose) addPart(data string, contentType string, contentEncoding string) {
	contentTypeSplit := strings.Split(contentType, "/")
	part := NewPartWithType(contentTypeSplit[0], contentTypeSplit[1])
	stream := NewMemStreamWithBuffer(data)
	//TODO: (Kane) should we close the stream?
	wrapper := NewDataWrapperWithStream(stream, contentEncoding)
	part.SetContentObject(wrapper)
	p.multipart.AddPart(part)

	if p.partsCount > 0 {
		// the message should be multipart
		p.message.SetMimePart(p.multipart)
	} else {
		// single part message
		p.message.SetMimePart(part)
	}
	p.partsCount++
}
