package gmime

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndMutationOnMime(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/inline-attachment.eml")
	assert.NoError(t, err)

	msg := Parse(string(mimeBytes))
	defer msg.Close()

	//Verify that we get subject and header parsed correctly
	assert.Equal(t, msg.Subject(), "test inline image attachment")
	assert.Equal(t, msg.ContentType(), "multipart/related")
	assert.Equal(t, msg.Header("Message-ID"), "<CAGPJ=uY91HEGoszHE9ELkB3wfcNJN4NGORM9q-vV8o_XJceBmg@mail.gmail.com>")

	contentType := []string{
		"text/plain",
		"text/html",
		"image/jpeg",
	}

	partText := []string{
		"kien image below\n\n[image: Inline image 1]\n\n--\nKien Pham\nSoftware Engineer, SendGrid\n",
		"<div dir=\"ltr\">kien image below<div><br></div><div><img src=\"cid:ii_1463f6eb06c77530\" alt=\"Inline image 1\" width=\"64\" height=\"64\"><br clear=\"all\"><div><br></div>-- <br><div dir=\"ltr\"><div>Kien Pham</div><div>Software Engineer, SendGrid<br>\n</div></div>\n</div></div>\n",
	}

	//Verify that we get parts contentType and text parsed correctly
	i := 0
	msg.Walk(func(p *Part) {
		assert.Equal(t, p.ContentType(), contentType[i])
		if p.IsText() {
			assert.Equal(t, p.Text(), partText[i])
			p.SetText(fmt.Sprintf("my replaced всякий текст スラングまで幅広く収録 (%d)", i))
		}
		i++
	})

	// Mutate subject header and body
	newSubject := "new subject"
	msg.SetSubject(newSubject)
	newMsgID := "new messageid"
	msg.SetHeader("Message-ID", newMsgID)

	// Verify subject/header and body are updated
	assert.Equal(t, msg.Subject(), newSubject)
	assert.Equal(t, msg.Header("Message-ID"), newMsgID)

	i = 0
	msg.Walk(func(p *Part) {
		if p.IsText() {
			assert.Equal(t, p.Text(), fmt.Sprintf("my replaced всякий текст スラングまで幅広く収録 (%d)", i))
		}
		i++
	})

	mimeBytes, err = msg.Export()
	assert.NoError(t, err)
	expectedMimeBytes, err := ioutil.ReadFile("test_data/inline-attachment_mutated.eml")
	assert.NoError(t, err)
	assert.Equal(t, len(mimeBytes), len(expectedMimeBytes))
}
