package gmime

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndMutationOnMime_Multipart(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/inline-attachment_multipart.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)
	defer msg.Close()

	//Verify that we get subject and header parsed correctly
	assert.Equal(t, msg.Subject(), "test inline image attachment")
	assert.Equal(t, msg.ContentType(), "multipart/alternative")
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
	var i, k int
	err = msg.Walk(func(p *Part) error {
		assert.Equal(t, contentType[i], p.ContentType())
		if p.IsText() {
			assert.Equal(t, partText[k], p.Text())
			p.SetText(fmt.Sprintf("my replaced всякий текст スラングまで幅広く収録 (%d)", i))
			k++
		}
		i++
		return nil
	})
	assert.NoError(t, err)

	msg.Walk(func(p *Part) error {
		if p.IsAttachment() {
			ct := p.ContentType()
			filename := p.Filename()
			assert.Equal(t, ct, "image/jpeg")
			assert.NotEqual(t, ct, "text/html")
			assert.NotEqual(t, ct, "text/plain")
			assert.Equal(t, "kien.jpg", filename)
		}

		return nil
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
	err = msg.Walk(func(p *Part) error {
		if p.IsText() {
			assert.Equal(t, p.Text(), fmt.Sprintf("my replaced всякий текст スラングまで幅広く収録 (%d)", i))
		}
		i++
		return nil
	})
	assert.NoError(t, err)
}

func TestParseAndMutationOnMime_NestedMultipart(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/inline-attachment_nested_multipart.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)
	defer msg.Close()

	//Verify that we get subject and header parsed correctly
	assert.Equal(t, msg.Subject(), "test inline image attachment")
	assert.Equal(t, msg.ContentType(), "multipart/related")
	assert.Equal(t, msg.Header("Message-ID"), "<CAGPJ=uY91HEGoszHE9ELkB3wfcNJN4NGORM9q-vV8o_XJceBmg@mail.gmail.com>")

	contentType := []string{
		"multipart/alternative",
		"text/plain",
		"text/html",
		"image/jpeg",
	}

	partText := []string{
		"kien image below\n\n[image: Inline image 1]\n\n--\nKien Pham\nSoftware Engineer, SendGrid\n",
		"<div dir=\"ltr\">kien image below<div><br></div><div><img src=\"cid:ii_1463f6eb06c77530\" alt=\"Inline image 1\" width=\"64\" height=\"64\"><br clear=\"all\"><div><br></div>-- <br><div dir=\"ltr\"><div>Kien Pham</div><div>Software Engineer, SendGrid<br>\n</div></div>\n</div></div>\n",
	}

	//Verify that we get parts contentType and text parsed correctly
	var i, k int
	err = msg.Walk(func(p *Part) error {
		assert.Equal(t, contentType[i], p.ContentType())
		if p.IsText() {
			assert.Equal(t, partText[k], p.Text())
			p.SetText(fmt.Sprintf("my replaced всякий текст スラングまで幅広く収録 (%d)", i))
			k++
		}
		i++
		return nil
	})
	assert.NoError(t, err)

	// Mutate subject header and body
	newSubject := "new subject"
	msg.SetSubject(newSubject)
	newMsgID := "new messageid"
	msg.SetHeader("Message-ID", newMsgID)

	// Verify subject/header and body are updated
	assert.Equal(t, msg.Subject(), newSubject)
	assert.Equal(t, msg.Header("Message-ID"), newMsgID)

	i = 0
	err = msg.Walk(func(p *Part) error {
		if p.IsText() {
			assert.Equal(t, p.Text(), fmt.Sprintf("my replaced всякий текст スラングまで幅広く収録 (%d)", i))
		}
		i++
		return nil
	})
	assert.NoError(t, err)
}

func TestAddHTMLAlternativeToPlainText(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/textplain.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)

	htmlPayload := "<html><body></body></html>"
	added := msg.AddHTMLAlternativeToPlainText(htmlPayload)
	assert.Equal(t, "multipart/alternative", msg.ContentType())
	assert.True(t, added)
	exported, err := msg.Export()
	assert.NoError(t, err)
	assert.Contains(t, string(exported), htmlPayload)
	msg.Close()

	mimeBytes, err = ioutil.ReadFile("test_data/inline-attachment_multipart.eml")
	assert.NoError(t, err)
	msg, err = Parse(string(mimeBytes))
	assert.NoError(t, err)
	added = msg.AddHTMLAlternativeToPlainText(htmlPayload)
	assert.Equal(t, "multipart/alternative", msg.ContentType())
	assert.False(t, added)
	exported, err = msg.Export()
	assert.NoError(t, err)
	assert.NotContains(t, string(exported), htmlPayload)
	msg.Close()
}

func TestRemoveAll(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)

	removed := msg.RemoveAllHeaders("X-HEADER")
	assert.Equal(t, "", msg.Header("X-HEADER"))
	assert.True(t, removed)

	removed = msg.RemoveAllHeaders("X-HEADER")
	assert.False(t, removed)

	assert.Equal(t, "Kien Pham <kien@sendgrid.com>", msg.Header("To"))
}

func TestAddAddresses(t *testing.T) {
	tests := []struct {
		header        string
		phrase        string
		address       string
		expectedError string
	}{
		{"to", "123", "to@to.com", ""},
		{"cc", "456", "cc@cc.com", ""},
		{"bcc", "789", "cc@cc.com", ""},
		{"from", "2342789", "from@from.com", ""},
		{"sender", "78119", "sender@sender.com", ""},
		{"reply-to", "734389", "reply-to@reply-to.com", ""},
		{"wtf", "999", "wtf@wtf.com", "unknown header wtf"},
	}

	for _, test := range tests {
		mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
		assert.NoError(t, err)
		msg, err := Parse(string(mimeBytes))
		assert.NoError(t, err)

		err = msg.AddAddress(test.header, test.phrase, test.address)
		if test.expectedError == "" {
			assert.NoError(t, err)

			to := msg.Header(test.header)
			assert.Contains(t, to, test.address)

			newMime, err := msg.Export()
			m := string(newMime)
			assert.NoError(t, err)
			assert.Contains(t, m, test.address)
			assert.Contains(t, m, test.phrase)
		} else {
			assert.Contains(t, err.Error(), test.expectedError)
		}
	}
}

func TestClearAddress(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)

	err = msg.ClearAddress("from")
	assert.NoError(t, err)
	err = msg.ClearAddress("to")
	assert.NoError(t, err)
	err = msg.ClearAddress("sender")
	assert.NoError(t, err)
	err = msg.ClearAddress("reply-to")
	assert.NoError(t, err)
	err = msg.ClearAddress("bcc")
	assert.NoError(t, err)
	err = msg.ClearAddress("cc")
	assert.NoError(t, err)
	err = msg.ClearAddress("wtf")
	assert.Contains(t, err.Error(), "unknown header wtf")

	newMime, err := msg.Export()
	m := string(newMime)
	assert.NotContains(t, m, "kien@sendgrid.com")
	assert.NotContains(t, m, "kpham@sendgrid.com")
	assert.NotContains(t, m, "kane@sendgrid.com")
	assert.NotContains(t, m, "isaac@sendgrid.com")
	assert.NotContains(t, m, "tim@sendgrid.com")
	assert.NotContains(t, m, "trevor@sendgrid.com")
}

func TestSetHeaderAddress(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)

	err = msg.SetHeader("from", "someone@somewhere.com")
	assert.Error(t, err)
	err = msg.SetHeader("sender", "someone@somewhere.com")
	assert.Error(t, err)
	err = msg.SetHeader("reply-to", "someone@somewhere.com")
	assert.Error(t, err)
	err = msg.SetHeader("to", "someone@somewhere.com")
	assert.Error(t, err)
	err = msg.SetHeader("cc", "someone@somewhere.com")
	assert.Error(t, err)
	err = msg.SetHeader("bcc", "someone@somewhere.com")
	assert.Error(t, err)
}
