package gmime

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndMutationOnMime_Multipart(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/inline-attachment_multipart.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)
	defer msg.Close()

	// Verify that we get subject and header parsed correctly
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

	// Verify that we get parts contentType and text parsed correctly
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

	// Verify that we get subject and header parsed correctly
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

	// Verify that we get parts contentType and text parsed correctly
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

func TestHeaders(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)

	assert.Equal(t, "<CAGPJ=uY91HEGoszHE9ELkB3wfcNJN4NGORM9q-vV8o_XJceBmg@mail.gmail.com>", msg.Headers().Get("Message-ID"))
	assert.Equal(t, "<CAGPJ=uY91HEGoszHE9ELkB3wfcNJN4NGORM9q-vV8o_XJceBmg@mail.gmail.com>", msg.Headers().Get("Message-Id"))

	headers := []string{}
	for name, values := range msg.Headers() {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s: %s", name, value))
		}
	}
	sort.Strings(headers)
	assert.Equal(t, []string{
		"Bcc: Tim <tim@sendgrid.com>",
		"Cc: trevor <trevor@sendgrid.com>",
		"Date: Tue, 27 May 2014 13:46:04 -0700",
		"Delivered-To: kien.pham@sendgrid.com",
		"From: Kien Pham <kpham@sendgrid.com>",
		"Message-Id: <CAGPJ=uY91HEGoszHE9ELkB3wfcNJN4NGORM9q-vV8o_XJceBmg@mail.gmail.com>",
		"Mime-Version: 1.0",
		"Received: by 10.220.115.205 with HTTP; Tue, 27 May 2014 13:46:04 -0700 (PDT)",
		"Reply-To: isaac <isaac@sendgrid.com>",
		"Sender: Kane Kim <kane@sendgrid.com>",
		"Subject: test inline image attachment",
		"To: Kien Pham <kien@sendgrid.com>",
		"X-Header: 1",
		"X-Header: 2",
		"X-Header: 3",
		"X-Header: 4",
	}, headers)
}

func TestRawHeaders(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)

	assert.Equal(t, " <CAGPJ=uY91HEGoszHE9ELkB3wfcNJN4NGORM9q-vV8o_XJceBmg@mail.gmail.com>\n", msg.RawHeaders().Get("Message-ID"))
	assert.Equal(t, " <CAGPJ=uY91HEGoszHE9ELkB3wfcNJN4NGORM9q-vV8o_XJceBmg@mail.gmail.com>\n", msg.RawHeaders().Get("Message-Id"))

	headers := []string{}
	for name, values := range msg.RawHeaders() {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%s:%s", name, value))
		}
	}
	sort.Strings(headers)
	assert.Equal(t, []string{
		"Bcc: Tim <tim@sendgrid.com>\n",
		"Cc: trevor <trevor@sendgrid.com>\n",
		"Date: Tue, 27 May 2014 13:46:04 -0700\n",
		"Delivered-To: kien.pham@sendgrid.com\n",
		"From: Kien Pham <kpham@sendgrid.com>\n",
		"Message-Id: <CAGPJ=uY91HEGoszHE9ELkB3wfcNJN4NGORM9q-vV8o_XJceBmg@mail.gmail.com>\n",
		"Mime-Version: 1.0\n",
		"Received: by 10.220.115.205 with HTTP; Tue, 27 May 2014 13:46:04 -0700 (PDT)\n",
		"Reply-To: isaac <isaac@sendgrid.com>\n",
		"Sender: Kane Kim <kane@sendgrid.com>\n",
		"Subject: test inline image attachment\n",
		"To: Kien Pham <kien@sendgrid.com>\n",
		"X-Header: 1\n",
		"X-Header: 2\n",
		"X-Header: 3\n",
		"X-Header: 4\n",
	}, headers)
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

func TestReplaceHeader(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)

	oldHeaders, err := headersSlice(mimeBytes)
	assert.NoError(t, err)

	// Replace the second X-HEADER with 5
	replace := "5"
	err = msg.ReplaceHeader("X-HEADER", "2", replace)
	assert.NoError(t, err)
	oldHeaders[13] = headerData{"X-HEADER", replace}
	mimeBytes, err = msg.Export()
	assert.NoError(t, err)
	newHeaders, err := headersSlice(mimeBytes)
	assert.NoError(t, err)
	// check order and value
	assert.True(t, equal(oldHeaders, newHeaders))

	// Replace a X-HEADER with value not exist in mime
	// Should fail, headers should not be changed
	err = msg.ReplaceHeader("X-HEADER", "value don't exist", replace)
	assert.Error(t, err)
	mimeBytes, err = msg.Export()
	assert.NoError(t, err)
	newHeaders, err = headersSlice(mimeBytes)
	assert.NoError(t, err)
	assert.True(t, equal(oldHeaders, newHeaders))

	// Replace a header that doesn't exist in mime
	// Should fail, headers should not be changed
	err = msg.ReplaceHeader("key don't exist", "1", replace)
	assert.Error(t, err)
	mimeBytes, err = msg.Export()
	assert.NoError(t, err)
	newHeaders, err = headersSlice(mimeBytes)
	assert.NoError(t, err)
	assert.True(t, equal(oldHeaders, newHeaders))
}

type headerData struct {
	name  string
	value string
}

func headersSlice(mimeBytes []byte) ([]headerData, error) {
	var headers []headerData
	b := bytes.NewReader(mimeBytes)
	scanner := bufio.NewScanner(b)

	for scanner.Scan() {
		header := strings.SplitN(scanner.Text(), ": ", 2)
		if len(header) == 2 {
			headers = append(headers, headerData{header[0], header[1]})
		}
	}
	return headers, scanner.Err()
}

func equal(a, b []headerData) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
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
		{"wtf", "999", "wtf@wtf.com", "can't add to header wtf"},
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

func TestParseAndAppendAddresses(t *testing.T) {
	tests := []struct {
		addresses string
		expected  string
	}{
		{"a@a.com", "a@a.com"},
		{"a@a.com,b@b.com", "a@a.com, b@b.com"},
		{"a@a.com b@b.com", "a@a.com, b@b.com"},
		{"a <a@a.com> b b@b.com", "a <a@a.com>"},
		{`a a@a.com, b <b@b.com>, "c" <c@c.com>`, "b <b@b.com>, c <c@c.com>"},
		{`a@a.com,b <b@b.com>`, "a@a.com, b <b@b.com>"},
		{`a@a.com,[] <badbrackets@b.com>, c <c@c.com>`, "a@a.com, c <c@c.com>"},
		{`a@a.com, "[]" <goodbrackets@b.com>, c@c.com`, `a@a.com, "[]" <goodbrackets@b.com>, c@c.com`},
	}

	for _, test := range tests {
		mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
		assert.NoError(t, err)
		msg, err := Parse(string(mimeBytes))
		assert.NoError(t, err)

		msg.RemoveHeader("to")
		msg.ParseAndAppendAddresses("to", test.addresses)
		assert.Equal(t, test.expected, msg.Header("to"))
	}
}

func TestIsAttachment(t *testing.T) {
	tests := []struct {
		filename     string
		isAttachment bool
	}{
		{"textplain.eml", false},
		{"multipleHeaders.eml", false},
		{"attachmentwithname.eml", true},
		{"attachmentwithoutname.eml", true},
		{"inlineattachment.eml", true},
		{"inline.eml", false},
	}

	for _, test := range tests {
		mimeBytes, err := ioutil.ReadFile(fmt.Sprintf("test_data/%s", test.filename))
		assert.NoError(t, err)
		msg, err := Parse(string(mimeBytes))
		assert.NoError(t, err)
		msg.Walk(func(p *Part) error {
			assert.Equal(t, test.isAttachment, p.IsAttachment())
			return nil
		})
	}
}

func TestIsLegacyAttachment(t *testing.T) {
	tests := []struct {
		filename     string
		contentType  string
		isAttachment bool
	}{
		{"nested_calendar.eml", "text/calendar", true},

		//include previous tests for backwards compatibility
		{"textplain.eml", "text/plain", false},
		{"multipleHeaders.eml", "text/plain", false},
		{"attachmentwithname.eml", "text/plain", true},
		{"attachmentwithoutname.eml", "text/plain", true},
		{"inlineattachment.eml", "text/plain", true},
		{"inline.eml", "text/plain", false},
	}

	for _, test := range tests {
		mimeBytes, err := ioutil.ReadFile(fmt.Sprintf("test_data/%s", test.filename))
		assert.NoError(t, err)
		msg, err := Parse(string(mimeBytes))
		assert.NoError(t, err)

		msg.WalkWithParent(func(p *Part) error {
			if p.ContentType() == test.contentType {
				assert.Equal(t, test.isAttachment, p.IsLegacyAttachment())
			}
			return nil
		})
	}
}

func TestParseAddressList(t *testing.T) {
	tests := []struct {
		addrList  string
		gAddrList []*mail.Address
	}{
		{
			addrList: "Foo Bar <foo@bar.baz>",
			gAddrList: []*mail.Address{
				{
					Name:    "Foo Bar",
					Address: "foo@bar.baz",
				},
			},
		},
		{
			addrList: "Foo Bar <foo@bar.baz>, Bar Baz <bar@foo.com>",
			gAddrList: []*mail.Address{
				{
					Name:    "Foo Bar",
					Address: "foo@bar.baz",
				},
				{
					Name:    "Bar Baz",
					Address: "bar@foo.com",
				},
			},
		},
		{
			addrList: "Foo Bar <foo@bar.baz>, Bar Baz <bar@foo.com>, Not an email at all",
			gAddrList: []*mail.Address{
				{
					Name:    "Foo Bar",
					Address: "foo@bar.baz",
				},
				{
					Name:    "Bar Baz",
					Address: "bar@foo.com",
				},
			},
		},
		{
			addrList: "Foo Bar <foo@bar.baz>, Bar Baz <bar@foo.com>, Another Email <another.email@mail.com>",
			gAddrList: []*mail.Address{
				{
					Name:    "Foo Bar",
					Address: "foo@bar.baz",
				},
				{
					Name:    "Bar Baz",
					Address: "bar@foo.com",
				},
				{
					Name:    "Another Email",
					Address: "another.email@mail.com",
				},
			},
		},
		{
			addrList: "<foo@bar.baz>, <bar@foo.baz>",
			gAddrList: []*mail.Address{
				{
					Address: "foo@bar.baz",
				},
				{
					Address: "bar@foo.baz",
				},
			},
		},
		{
			addrList: "foo@bar.baz, <bar@foo.baz>",
			gAddrList: []*mail.Address{
				{
					Address: "foo@bar.baz",
				},
				{
					Address: "bar@foo.baz",
				},
			},
		},
		{
			addrList: "foo@bar.baz, Bar Foo <bar@foo.baz>",
			gAddrList: []*mail.Address{
				{
					Address: "foo@bar.baz",
				},
				{
					Name:    "Bar Foo",
					Address: "bar@foo.baz",
				},
			},
		},
		{
			addrList: "foo@bar.baz, Bar Foo bar@foo.baz",
			gAddrList: []*mail.Address{
				{
					Address: "foo@bar.baz",
				},
			},
		},
		{
			addrList: "foo@bar.baz, bar@foo.baz",
			gAddrList: []*mail.Address{
				{
					Address: "foo@bar.baz",
				},
				{
					Address: "bar@foo.baz",
				},
			},
		},
	}

	for _, test := range tests {
		got := ParseAddressList(test.addrList)
		assert.Equal(t, test.gAddrList, got)
	}
}

func TestAppendAddressList(t *testing.T) {
	tests := []struct {
		addrs  []*mail.Address
		header string
	}{
		{
			header: "Foo Bar <foo@bar.baz>, Bar Baz <bar@foo.com>, Another Email\t<another.email@mail.com>",
			addrs: []*mail.Address{
				{
					Name:    "Foo Bar",
					Address: "foo@bar.baz",
				},
				{
					Name:    "Bar Baz",
					Address: "bar@foo.com",
				},
				{
					Name:    "Another Email",
					Address: "another.email@mail.com",
				},
			},
		},
		{
			header: "Foo Bar <foo@bar.baz>, Bar Baz <bar@foo.com>",
			addrs: []*mail.Address{
				{
					Name:    "Foo Bar",
					Address: "foo@bar.baz",
				},
				{
					Name:    "Bar Baz",
					Address: "bar@foo.com",
				},
			},
		},
		// This is an actual test, no addrs == empty header
		{},
	}

	for _, test := range tests {
		mimeBytes, err := ioutil.ReadFile("test_data/multipleHeaders.eml")
		assert.NoError(t, err)
		msg, err := Parse(string(mimeBytes))
		assert.NoError(t, err)

		msg.RemoveHeader("from")
		err = msg.AppendAddressList("from", test.addrs)
		assert.NoError(t, err)

		assert.Equal(t, test.header, msg.Header("from"))
	}
}

// TestPart_ContentID not only tests the function fetches content-id correctly,
// but it also makes sure there is no crash when the envelope is closed (since we have seen a bug around this)
func TestPart_ContentID(t *testing.T) {
	mimeBytes, err := ioutil.ReadFile(`test_data/attachment-content-id.eml`)
	assert.NoError(t, err)
	msg, err := Parse(string(mimeBytes))
	assert.NoError(t, err)

	actual := ""
	expected := `f_l0eugokh0`
	err = msg.Walk(func(p *Part) error {
		if p.IsAttachment() {
			actual = p.ContentID()
		}
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	msg.Close()
}

func TestPart_String(t *testing.T) {
	type TestCase struct {
		filename string
		expected map[string]string
	}
	tcs := []TestCase{
		{
			filename: `test_data/inline-attachment_multipart.eml`,
			expected: map[string]string{
				"image/jpeg": "Content-Type: image/jpeg; name=\"kien.jpg\"\r\nContent-Disposition: attachment; filename=\"kien.jpg\"\r\nContent-Transfer-Encoding: base64\r\nContent-ID: <ii_1463f6eb06c77530>\r\nX-Attachment-Id: ii_1463f6eb06c77530\r\n\r\n/9j/4AAQSkZJRgABAQAAAQABAAD/4QDQRXhpZgAASUkqAAgAAAADABIBAwABAAAAAQAAADEBAgAH\r\nAAAAMgAAAGmHBAABAAAAOgAAAAAAAABQaWNhc2EAAAYAAJAHAAQAAAAwMjIwAaADAAEAAAABAAAA\r\nAqAEAAEAAABAAAAAA6AEAAEAAABAAAAABaAEAAEAAACqAAAAIKQCACEAAACIAAAAAAAAADUyYmMz\r\nZTk3NjEyNzFlY2IwMDAwMDAwMDAwMDAwMDAwAAACAAEAAgAEAAAAUjk4AAIABwAEAAAAMDEwMAAA\r\nAAD/4QEiaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLwA8P3hwYWNrZXQgYmVnaW49Iu+7vyIg\r\naWQ9Ilc1TTBNcENlaGlIenJlU3pOVGN6a2M5ZCI/PiA8eDp4bXBtZXRhIHhtbG5zOng9ImFkb2Jl\r\nOm5zOm1ldGEvIiB4OnhtcHRrPSJYTVAgQ29yZSA1LjEuMiI+IDxyZGY6UkRGIHhtbG5zOnJkZj0i\r\naHR0cDovL3d3dy53My5vcmcvMTk5OS8wMi8yMi1yZGYtc3ludGF4LW5zIyI+IDxyZGY6RGVzY3Jp\r\ncHRpb24gcmRmOmFib3V0PSIiLz4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gICA8P3hwYWNrZXQg\r\nZW5kPSJ3Ij8+/9sAhAADAgIDAgIDAwMDBAMDBAUIBQUEBAUKBwcGEBUKFRULCgsLDRANDRQXDxcL\r\nDg0WDhEQDxUVDxAPDhgQDBgPEA8PAQMEBAYFBgoGBgoMDQwNDg8QDQ0PDQ0MDQwQDQ8NDQ0NDQ0N\r\nDQwMDAwMDA4MDAwNDA0NDQwNDA0NDA0MDAwPDAz/wAARCABAAEADAREAAhEBAxEB/8QAGwAAAgMB\r\nAQEAAAAAAAAAAAAABwgEBgkFAgP/xABGEAABAwIEAwQDCwcNAAAAAAACAQMEBRIABgcREyEiCDEy\r\nQiNBUgkUFRZDUWJxcoGiJDNhc4KSwhclNERjZHSToaOxs9L/xAAbAQABBQEBAAAAAAAAAAAAAAAE\r\nAgMFBgcBAP/EACkRAAEEAQMDAwQDAAAAAAAAAAEAAgMEEQUSIRMyMyIjMUFRsfAVQqH/2gAMAwEA\r\nAhEDEQA/AHJ1MMBz1WQL52/+scWugMxKqWm4lVKqlYYo8B2XLRWgbS5fn79v4sTHaxBEZeg1rd2s\r\n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\r\nesHhf//Z",
				"text/html":  "Content-Type: text/html; charset=UTF-8\r\n\r\n<div dir=\"ltr\">kien image below<div><br></div><div><img src=\"cid:ii_1463f6eb06c77530\" alt=\"Inline image 1\" width=\"64\" height=\"64\"><br clear=\"all\"><div><br></div>-- <br><div dir=\"ltr\"><div>Kien Pham</div><div>Software Engineer, SendGrid<br>\r\n</div></div>\r\n</div></div>",
				"text/plain": "Content-Type: text/plain; charset=UTF-8\r\n\r\nkien image below\r\n\r\n[image: Inline image 1]\r\n\r\n--\r\nKien Pham\r\nSoftware Engineer, SendGrid",
			},
		},
		{
			filename: `test_data/rfc822.eml`,
			expected: map[string]string{
				"message/feedback-report": "Content-Disposition: inline\r\nContent-Type: message/feedback-report\r\n\r\nFeedback-Type: abuse\r\nUser-Agent: AOL SComp\r\nVersion: 0.1\r\nReceived-Date: Thu,  2 Sep 2010 15:12:11 -0400 (EDT)\r\nSource-IP: 74.63.231.149\r\nReported-Domain: o7463231149.static.reverse.sendgrid.net\r\nRedacted-Address: redacted\r\nRedacted-Address: redacted@",
				"message/rfc822":          "Content-Type: message/rfc822\r\nContent-Disposition: inline\r\n\r\nReturn-Path: <bounces+15167-8e67-redacted=aol.com@sendgrid.me>\r\nReceived: from mtain-di06.r1000.mx.aol.com (mtain-di06.r1000.mx.aol.com [172.29.64.10]) by air-da07.mail.aol.com (v129.4) with ESMTP id MAILINDA073-863f4c7ff9f999; Thu, 02 Sep 2010 15:24:41 -0400\r\nReceived: from o7463231149.static.reverse.sendgrid.net (o7463231149.static.reverse.sendgrid.net [74.63.231.149])\r\n\tby mtain-di06.r1000.mx.aol.com (Internet Inbound) with SMTP id B984D3800C101\r\n\tfor <redacted>; Thu,  2 Sep 2010 15:12:11 -0400 (EDT)\r\nDKIM-Signature: v=1; a=rsa-sha1; c=relaxed; d=sendgrid.me; h=\r\n\tcontent-type:content-transfer-encoding:mime-version:from:to\r\n\t:subject:message-id:date:sender:list-unsubscribe; s=smtpapi; bh=\r\n\t2XXf6FQE79UDRHRZkGBgP/3VOPQ=; b=UaGTnfj/JowTYMqpJwXc87tk4BWt6kl4\r\n\tYvxz0HhScM+2xQtQTiDxPWOSFnLyTey6azzJysZw4mGZ0T5vwBcgCDe2ldzBUZVP\r\n\tB4E5EdpO4opU6/2FcYaMwdfwc96lliHAnYhqDHO2AAum6SxwWOHj5Z0+lYjVfQCL\r\n\tAhsV/Y2HgR8=\r\nDomainKey-Signature: a=rsa-sha1; c=nofws; d=sendgrid.me; h=content-type\r\n\t:content-transfer-encoding:mime-version:from:to:subject\r\n\t:message-id:date:sender:list-unsubscribe; q=dns; s=smtpapi; b=mI\r\n\tpOeRnkiGAJ8ksC87pgKrZM0CaHC486kj6kHuR/SNADBdqcWX2SWiGnTnNIDD3PGN\r\n\t/xlsNDzDDX085eCSDLFtb8qIqjTg04VYmwZCEzIHEsXYGVCeVwy+HelrjzaZWvH8\r\n\tPjPuIBU9SHos8mXlL+jFo1gkkVZydY8gjM6P6zWwI=\r\nReceived: by 10.36.109.145 with SMTP id mf12.5291.4C7FE63EA\r\n        Thu, 02 Sep 2010 11:00:30 -0700 (PDT)\r\nContent-Type: multipart/alternative;\r\n boundary=\"----------=_1283450430-3201-82\";\r\n charset=\"UTF-8\"\r\nMIME-Version: 1.0\r\nX-Mailer: MIME-tools 5.428 (Entity 5.428)\r\nFrom: Tickfolio <info@tickfolio.com>\r\nTo: redacted@aol.com\r\nSubject: Tickfolio.com - Great Deals on Bears,  Ohio State and Notre Dame\r\n Football Tickets - 25% off at Cookies by Joey $50 voucher for $37 .50\r\nMessage-ID: <1283450583.6569568232097351@mf12.sendgrid.net>\r\nDate: Thu, 02 Sep 2010 11:03:03 -0700 (PDT)\r\nX-Sendgrid-EID: XXdhm96AWnTxnrn3EJ5eiKBJHPkRkfr8a8V2FwfC2nt77Ld3IpoIo08eEU6nkzNz/D3d+bWm5/nw1V5YMb7f1H5xCFEtYTrUGQaoogv9tgQZpkNWZrZyNQmRNuIZjLC+7WeHCvUjnY8Qp82Ac1c1ew== \r\nX-Sendgrid-ID: VPWZYjw6GOzHdwkwPeoX9QiEbzQXX/gF9P8njHP5+LAmfPHZNU0Z1CE2j34aFcymAfzbSNgCfwklPF6sm+frHw5ZuQoyDD0xhdGt3g6VSDM=\r\nSender: Tickfolio <info=tickfolio.com@sendgrid.me>\r\nList-Unsubscribe: <http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/um>, <mailto:unsubscribe@sendgrid.me?subject=http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/um>\r\nx-aol-global-disposition: G\r\nX-AOL-SCOLL-AUTHENTICATION: mail_rly_antispam_dkim-d247.2 ; domain : sendgrid.me DKIM : pass\r\nx-aol-sid: 3039ac1d400a4c7ff70b73e2\r\nX-AOL-IP: 74.63.231.149\r\nX-AOL-SPF: domain : sendgrid.me SPF : pass\r\nContent-Transfer-Encoding: 7bit\r\n\r\n\n\r\n------------=_1283450430-3201-82\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\nContent-Disposition: inline\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nTickfolio.com - Hot Deals\r\n\r\n\r\n\r\n$50 Value =E2=80=93 for $37.50 -25% Discount =E2=80=93 You Save $12.50\r\n\r\n\r\n\r\nGourmet Cookies Shipped Nationwide. Best Wishes, Thank You, Congrats, and=\r\n of course Happy Birthday... There are endless reasons to send your friend=\r\ns, family members, co-workers and clients they are truly amazing!\r\n\r\n\r\n\r\n-=3D Just Great Tickets =3D-\r\n\r\n\r\n\r\nTickfolio in partnership with Just Great Tickets offers outstanding deals=\r\n on Bears, Sox, Cubs, Ohio State and Nortre Dame Football tickets, Bulls,=\r\n Blackhawks and ALL Concerts and Theater for Individuals and Groups! (Go=\r\n to http://tickfolio.com/archives/1106?utm_source=3Dsendgrid.com&utm_mediu=\r\nm=3Demail&utm_campaign=3Dwebsite to complete the inquiry form and receive=\r\n a quote within 24 hours)\r\n\r\n\r\n\r\n-=3D JunoWallet =3D-\r\n\r\n\r\n\r\nTickfolio in partnership with JunoWallet is pleased to offer FREE Promotio=\r\nnal Mobile Gift Cards delivered by Tickfolio.com.\r\n\r\n\r\n\r\nVisit tickfolio.com for more details.If you'd like to unsubscribe and stop=\r\n receiving these emails click here: http://u15167.sendgrid.org/s/SeMzKTgIS=\r\nMiH7cKVML-Qdg/ut.\r\n\r\n------------=_1283450430-3201-82\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nContent-Disposition: inline\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<title></title>\r\n\r\n\r\n\r\n<table width=3D\"100%\" bgcolor=3D\"#e7e9eb\" border=3D\"0\" cellpadding=3D\"0\"=\r\n cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" align=3D\"center\" bgcolor=3D\"#e7e9eb\" border=3D\"0\" bor=\r\ndercolor=3D\"#000000\" cellpadding=3D\"0\" cellspacing=3D\"0\"><tbody><tr><td va=\r\nlign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"0=\r\n\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"504\"><a href=3D=\r\n\"http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/h0\"><img  style=3D\"di=\r\nsplay: block;\" alt=3D\"\" src=3D\"https://app.icontact.com/icp/loadimage.php/=\r\nmogile/615954/d21f556e4ec08604ef1a4436502ad360/image/jpeg\" complete=3D\"com=\r\nplete\" width=3D\"504\" border=3D\"0\" height=3D\"90\"></a></td><!-- br\r\n\r\n--><td valign=3D\"top\" width=3D\"196\">\r\n\r\n<table width=3D\"196\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"0=\r\n\" cellspacing=3D\"0\"><tbody><tr>&nbsp;</tr><tr><td>\r\n\r\n<table width=3D\"196\" cellpadding=3D\"0\" cellspacing=3D\"0\"><tbody><tr><td va=\r\nlign=3D\"top\" width=3D\"100%\"><a href=3D\"\"><img  alt=3D\"\" src=3D\"https://app=\r\n.icontact.com/icp/loadimage.php/mogile/615954/34401e3f085814b4d8f32965c6a0=\r\n7562/image/jpeg\" complete=3D\"complete\" width=3D\"123\" border=3D\"0\" height=\r\n=3D\"46\"></a></td><td valign=3D\"top\" width=3D\"100%\"><a href=3D\"http://u1516=\r\n7.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/h1\"><img  alt=3D\"\" src=3D\"https://=\r\napp.icontact.com/icp/loadimage.php/mogile/615954/c2cacae59990552b6627a5cc8=\r\nd54edf8/image/jpeg\" complete=3D\"complete\" width=3D\"73\" border=3D\"0\" height=\r\n=3D\"46\"></a></td></tr></tbody></table></td></tr><tr><td>\r\n\r\n<table width=3D\"196\" cellpadding=3D\"10\" cellspacing=3D\"0\"><tbody><tr><td=\r\n align=3D\"center\"><a style=3D\"text-decoration: none;\" href=3D\"f2f_url\"><sp=\r\nan color=3D\"#333333\" style=3D\"font-family: Arial, Helvetica, sans-serif\"=\r\n size=3D\"2;\"><strong>Share With a Friend</strong></span></a> </td></tr></t=\r\nbody></table></td></tr></tbody></table></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr><tr><td valign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" bgcolor=3D\"#c7ced1\" border=3D\"0\" bordercolor=3D\"#0000=\r\n00\" cellpadding=3D\"0\" cellspacing=3D\"1\" height=3D\"40\"><tbody><tr><td valig=\r\nn=3D\"middle\" width=3D\"60\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"><a style=3D=\r\n\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/SeMzKTgISMiH=\r\n7cKVML-Qdg/h2\"><span style=3D\"font-size: 12px;\" color=3D\"#434343\" style=3D=\r\n\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Home</span></a></t=\r\nd><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"75\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h3\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Sales=\r\n Agent</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"65\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h4\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">About<=\r\n/span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"70\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h5\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Contac=\r\nt</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"55\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h6\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">FAQ</s=\r\npan></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"145\" align=3D\"center\" bgcolor=3D\"#f3f4f5=\r\n\"><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s=\r\n/SeMzKTgISMiH7cKVML-Qdg/h7\"><span style=3D\"font-size: 12px;\" color=3D\"#434=\r\n343\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Featu=\r\nre Your Business</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"80\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h8\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Partne=\r\nrships</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"130\" align=3D\"center\" bgcolor=3D\"#f3f4f5=\r\n\"><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s=\r\n/SeMzKTgISMiH7cKVML-Qdg/h9\"><span style=3D\"font-size: 12px;\" color=3D\"#434=\r\n343\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Sugge=\r\nst A Business</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"10\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><img  style=3D\"display: block;\" alt=3D\"\" src=3D\"https://app.icontact.com/=\r\nicp/loadimage.php/mogile/372789/fcaeb1c2d6caba6ffebdeb39c7149266/image/gif=\r\n\" complete=3D\"complete\" width=3D\"1\" height=3D\"38\"></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr><tr><td valign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"0=\r\n\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"7\"><img  style=\r\n=3D\"display: block;\" alt=3D\"\" src=3D\"https://app.icontact.com/icp/loadimag=\r\ne.php/mogile/372789/fcaeb1c2d6caba6ffebdeb39c7149266/image/gif\" complete=\r\n=3D\"complete\" width=3D\"1\" height=3D\"1\"></td><!-- br\r\n\r\n--><td style=3D\"padding: 20px;\" valign=3D\"middle\" width=3D\"653\"><font colo=\r\nr=3D\"#000000\" face=3D\"Arial, Helvetica, sans-serif\" size=3D\"4\"><strong>\r\n\r\n<h2><font size=3D\"4\">$50 Value =E2=80=93 for $37.50 -25% Discount =E2=80=\r\n=93 You Save $12.50</font></h2>\r\n\r\n<h3><font size=3D\"3\">Gourmet Cookies Shipped Nationwide. Best Wishes, Than=\r\nk You, Congrats, and of course Happy Birthday... There are endless reasons=\r\n to send your friends, family members, co-workers and clients they are tru=\r\nly amazing!</font></h3>\r\n\r\n<h3>About Cookies By Joey</h3>\r\n\r\n<div><font size=3D\"3\">Joanne Sherman, affectionately known as =E2=80=9CJoe=\r\ny,=E2=80=9D was raised in a large Italian family with 10 children. The kit=\r\nchen was a tremendous part of their lives, as Joey=E2=80=99s mother prepar=\r\ned enormous meals daily. A quick study in her mother=E2=80=99s kitchen, Jo=\r\ney followed her mother=E2=80=99s culinary prowess and turned her attention=\r\n to baking. She rapidly moved up the ranks in the bakery department of the=\r\n second largest supermarket chain in the Chicagoland area, and eventually=\r\n moved on to focus on her own creations. Over the years, Joey honed her sk=\r\nills and developed her own signature tastes that are celebrated in many ba=\r\nking goods=E2=80=A6most notably, Joey=E2=80=99s cookies.</font></div>\r\n\r\n<div>&nbsp;</div>\r\n\r\n<div>\r\n\r\n<h3>Testimonials</h3>\r\n\r\n<p><font size=3D\"3\">=E2=80=9CThank you for the most amazing basket of cook=\r\nies . . . they are absolutely =E2=80=9Cmelt in your mouth=E2=80=9D delicio=\r\nus!! Our family can=E2=80=99t get enough of them. Your beautiful packaging=\r\n is out of this world and your attention to detail is spectacular! We appr=\r\neciate all the care and lovely wrapping that came with the cookies. The qu=\r\nality of your product and your hands on customer service truly set you apa=\r\nrt from any other company. We look forward to sending your cookies to fami=\r\nly, friends and clients for all our special occasions and holidays. We are=\r\n customers for life!! Thank you again, and we look forward to eating more=\r\n Joey=E2=80=99s Cookies!=E2=80=9D</font></p></div></strong></font>\r\n\r\n<div><span color=3D\"#000000\" style=3D\"font-family: Arial, Helvetica, sans-=\r\nserif\" size=3D\"4;\"><strong>Go to&nbsp;</strong></span><strong><span  style=\r\n=3D\"font-family: Arial\" size=3D\"4;\">&nbsp;</span></strong></div>\r\n\r\n<div><strong><span  style=3D\"font-family: Arial\" size=3D\"4;\"><a href=3D\"ht=\r\ntp://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/h10\">www.tickfolio.com/a=\r\nrchives/1106</a> &nbsp;to&nbsp;order your Gourmet Cookies today!</span></s=\r\ntrong></div>\r\n\r\n<div><strong><span  style=3D\"font-family: Arial\" size=3D\"4;\">&nbsp;</span>=\r\n</strong></div>\r\n\r\n<div><strong><span  style=3D\"font-family: Arial;\">Comments or Questions?=\r\n <a href=3D\"mailto:info@Tickfolio.com\">info@Tickfolio.com</a></span></stro=\r\nng></div></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr><tr><td valign=3D\"top\" width=3D\"100%\"><img  style=3D\"display: bloc=\r\nk;\" alt=3D\"\" src=3D\"https://app.icontact.com/icp/loadimage.php/mogile/6159=\r\n54/b4bc5d5af51f975c01c81e2a038f6a25/image/jpeg\" complete=3D\"complete\"></td=\r\n><!-- br\r\n\r\n--></tr><tr><td valign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"0=\r\n\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"69\"><img  style=\r\n=3D\"display: block;\" alt=3D\"\" src=3D\"https://app.icontact.com/icp/loadimag=\r\ne.php/mogile/615954/fba9b691859d6cdfb483e09e57fda789/image/jpeg\" complete=\r\n=3D\"complete\" width=3D\"69\" height=3D\"279\"></td><!-- br\r\n\r\n--><td valign=3D\"top\" width=3D\"278\">\r\n\r\n<table width=3D\"278\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"2=\r\n0\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"100%\"><span st=\r\nyle=3D\"font-size: 28px;\" color=3D\"#000000\" style=3D\"font-family: Georgia,=\r\n Times New Roman, Times, serif\" size=3D\"5;\">Deal of the Week</span><br>\r\n\r\n<font color=3D\"#000000\" face=3D\"Arial, Helvetica, sans-serif\" size=3D\"2\">T=\r\nickfolio in partnership with Just Great Tickets offers outstanding deals=\r\n on Bears, Sox, Cubs, Ohio State and Nortre Dame Football tickets, Bulls,=\r\n Blackhawks and ALL Concerts and Theater for Individuals and Groups! (<str=\r\nong><a href=3D\"http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/h11\">Cl=\r\nick here</a></strong> to complete the inquiry form and receive a quote wit=\r\nhin 24 hours)<br>\r\n\r\n<br>\r\n\r\n<a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/S=\r\neMzKTgISMiH7cKVML-Qdg/h12\"><img  style=3D\"display: block;\" alt=3D\"\" src=3D=\r\n\"https://app.icontact.com/icp/loadimage.php/mogile/615954/27cade05fad83fd6=\r\n736d3907aab81341/image/jpeg\" complete=3D\"complete\" align=3D\"right\" border=\r\n=3D\"0\"></a></font></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--><td valign=3D\"top\" width=3D\"66\"><img  style=3D\"display: block;\" alt=3D\"=\r\n\" src=3D\"https://app.icontact.com/icp/loadimage.php/mogile/615954/fbf9f554=\r\nc74eb4d6b351237f5cfb2a6e/image/jpeg\" complete=3D\"complete\" width=3D\"66\" he=\r\night=3D\"286\"></td><!-- br\r\n\r\n--><td valign=3D\"top\" width=3D\"287\">\r\n\r\n<table width=3D\"287\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"2=\r\n0\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"100%\"><span st=\r\nyle=3D\"font-size: 28px;\" color=3D\"#000000\" style=3D\"font-family: Georgia,=\r\n Times New Roman, Times, serif\" size=3D\"5;\">JunoWallet</span><br>\r\n\r\n<font color=3D\"#000000\" face=3D\"Arial, Helvetica, sans-serif\" size=3D\"2\">T=\r\nickfolio in partnership with JunoWallet is pleased to offer FREE Promotion=\r\nal Mobile Gift Cards delivered by <strong><a href=3D\"http://u15167.sendgri=\r\nd.org/s/SeMzKTgISMiH7cKVML-Qdg/h13\">Tickfolio.com</a></strong>. No purchas=\r\ne is required to receive your free gift cards so enjoy our generous gift=\r\n to you our valued clients and prospects...<br>\r\n\r\n<br>\r\n\r\n<a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/S=\r\neMzKTgISMiH7cKVML-Qdg/h14\"><img  style=3D\"display: block;\" alt=3D\"\" src=3D=\r\n\"https://app.icontact.com/icp/loadimage.php/mogile/615954/27cade05fad83fd6=\r\n736d3907aab81341/image/jpeg\" complete=3D\"complete\" align=3D\"right\" border=\r\n=3D\"0\"></a></font></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr></tbody></table><br>\r\n\r\n<span style=3D\"padding: 0px;\">&nbsp;</span></td><!-- br\r\n\r\n--></tr></tbody></table>If you'd like to unsubscribe and stop receiving th=\r\nese emails <a href=3D\"http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/=\r\nuh\">click here</a>.\r\n\r\n<img src=3D\"http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/o0.gif\">\r\n\r\n------------=_1283450430-3201-82--",
				"text/html":               "Content-Type: text/html; charset=\"UTF-8\"\r\nContent-Disposition: inline\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n<title></title>\r\n\r\n\r\n\r\n<table width=3D\"100%\" bgcolor=3D\"#e7e9eb\" border=3D\"0\" cellpadding=3D\"0\"=\r\n cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" align=3D\"center\" bgcolor=3D\"#e7e9eb\" border=3D\"0\" bor=\r\ndercolor=3D\"#000000\" cellpadding=3D\"0\" cellspacing=3D\"0\"><tbody><tr><td va=\r\nlign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"0=\r\n\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"504\"><a href=3D=\r\n\"http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/h0\"><img  style=3D\"di=\r\nsplay: block;\" alt=3D\"\" src=3D\"https://app.icontact.com/icp/loadimage.php/=\r\nmogile/615954/d21f556e4ec08604ef1a4436502ad360/image/jpeg\" complete=3D\"com=\r\nplete\" width=3D\"504\" border=3D\"0\" height=3D\"90\"></a></td><!-- br\r\n\r\n--><td valign=3D\"top\" width=3D\"196\">\r\n\r\n<table width=3D\"196\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"0=\r\n\" cellspacing=3D\"0\"><tbody><tr>&nbsp;</tr><tr><td>\r\n\r\n<table width=3D\"196\" cellpadding=3D\"0\" cellspacing=3D\"0\"><tbody><tr><td va=\r\nlign=3D\"top\" width=3D\"100%\"><a href=3D\"\"><img  alt=3D\"\" src=3D\"https://app=\r\n.icontact.com/icp/loadimage.php/mogile/615954/34401e3f085814b4d8f32965c6a0=\r\n7562/image/jpeg\" complete=3D\"complete\" width=3D\"123\" border=3D\"0\" height=\r\n=3D\"46\"></a></td><td valign=3D\"top\" width=3D\"100%\"><a href=3D\"http://u1516=\r\n7.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/h1\"><img  alt=3D\"\" src=3D\"https://=\r\napp.icontact.com/icp/loadimage.php/mogile/615954/c2cacae59990552b6627a5cc8=\r\nd54edf8/image/jpeg\" complete=3D\"complete\" width=3D\"73\" border=3D\"0\" height=\r\n=3D\"46\"></a></td></tr></tbody></table></td></tr><tr><td>\r\n\r\n<table width=3D\"196\" cellpadding=3D\"10\" cellspacing=3D\"0\"><tbody><tr><td=\r\n align=3D\"center\"><a style=3D\"text-decoration: none;\" href=3D\"f2f_url\"><sp=\r\nan color=3D\"#333333\" style=3D\"font-family: Arial, Helvetica, sans-serif\"=\r\n size=3D\"2;\"><strong>Share With a Friend</strong></span></a> </td></tr></t=\r\nbody></table></td></tr></tbody></table></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr><tr><td valign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" bgcolor=3D\"#c7ced1\" border=3D\"0\" bordercolor=3D\"#0000=\r\n00\" cellpadding=3D\"0\" cellspacing=3D\"1\" height=3D\"40\"><tbody><tr><td valig=\r\nn=3D\"middle\" width=3D\"60\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"><a style=3D=\r\n\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/SeMzKTgISMiH=\r\n7cKVML-Qdg/h2\"><span style=3D\"font-size: 12px;\" color=3D\"#434343\" style=3D=\r\n\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Home</span></a></t=\r\nd><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"75\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h3\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Sales=\r\n Agent</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"65\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h4\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">About<=\r\n/span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"70\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h5\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Contac=\r\nt</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"55\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h6\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">FAQ</s=\r\npan></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"145\" align=3D\"center\" bgcolor=3D\"#f3f4f5=\r\n\"><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s=\r\n/SeMzKTgISMiH7cKVML-Qdg/h7\"><span style=3D\"font-size: 12px;\" color=3D\"#434=\r\n343\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Featu=\r\nre Your Business</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"80\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/=\r\nSeMzKTgISMiH7cKVML-Qdg/h8\"><span style=3D\"font-size: 12px;\" color=3D\"#4343=\r\n43\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Partne=\r\nrships</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"130\" align=3D\"center\" bgcolor=3D\"#f3f4f5=\r\n\"><a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s=\r\n/SeMzKTgISMiH7cKVML-Qdg/h9\"><span style=3D\"font-size: 12px;\" color=3D\"#434=\r\n343\" style=3D\"font-family: Arial, Helvetica, sans-serif\" size=3D\"2;\">Sugge=\r\nst A Business</span></a></td><!-- br\r\n\r\n--><td valign=3D\"middle\" width=3D\"10\" align=3D\"center\" bgcolor=3D\"#f3f4f5\"=\r\n><img  style=3D\"display: block;\" alt=3D\"\" src=3D\"https://app.icontact.com/=\r\nicp/loadimage.php/mogile/372789/fcaeb1c2d6caba6ffebdeb39c7149266/image/gif=\r\n\" complete=3D\"complete\" width=3D\"1\" height=3D\"38\"></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr><tr><td valign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"0=\r\n\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"7\"><img  style=\r\n=3D\"display: block;\" alt=3D\"\" src=3D\"https://app.icontact.com/icp/loadimag=\r\ne.php/mogile/372789/fcaeb1c2d6caba6ffebdeb39c7149266/image/gif\" complete=\r\n=3D\"complete\" width=3D\"1\" height=3D\"1\"></td><!-- br\r\n\r\n--><td style=3D\"padding: 20px;\" valign=3D\"middle\" width=3D\"653\"><font colo=\r\nr=3D\"#000000\" face=3D\"Arial, Helvetica, sans-serif\" size=3D\"4\"><strong>\r\n\r\n<h2><font size=3D\"4\">$50 Value =E2=80=93 for $37.50 -25% Discount =E2=80=\r\n=93 You Save $12.50</font></h2>\r\n\r\n<h3><font size=3D\"3\">Gourmet Cookies Shipped Nationwide. Best Wishes, Than=\r\nk You, Congrats, and of course Happy Birthday... There are endless reasons=\r\n to send your friends, family members, co-workers and clients they are tru=\r\nly amazing!</font></h3>\r\n\r\n<h3>About Cookies By Joey</h3>\r\n\r\n<div><font size=3D\"3\">Joanne Sherman, affectionately known as =E2=80=9CJoe=\r\ny,=E2=80=9D was raised in a large Italian family with 10 children. The kit=\r\nchen was a tremendous part of their lives, as Joey=E2=80=99s mother prepar=\r\ned enormous meals daily. A quick study in her mother=E2=80=99s kitchen, Jo=\r\ney followed her mother=E2=80=99s culinary prowess and turned her attention=\r\n to baking. She rapidly moved up the ranks in the bakery department of the=\r\n second largest supermarket chain in the Chicagoland area, and eventually=\r\n moved on to focus on her own creations. Over the years, Joey honed her sk=\r\nills and developed her own signature tastes that are celebrated in many ba=\r\nking goods=E2=80=A6most notably, Joey=E2=80=99s cookies.</font></div>\r\n\r\n<div>&nbsp;</div>\r\n\r\n<div>\r\n\r\n<h3>Testimonials</h3>\r\n\r\n<p><font size=3D\"3\">=E2=80=9CThank you for the most amazing basket of cook=\r\nies . . . they are absolutely =E2=80=9Cmelt in your mouth=E2=80=9D delicio=\r\nus!! Our family can=E2=80=99t get enough of them. Your beautiful packaging=\r\n is out of this world and your attention to detail is spectacular! We appr=\r\neciate all the care and lovely wrapping that came with the cookies. The qu=\r\nality of your product and your hands on customer service truly set you apa=\r\nrt from any other company. We look forward to sending your cookies to fami=\r\nly, friends and clients for all our special occasions and holidays. We are=\r\n customers for life!! Thank you again, and we look forward to eating more=\r\n Joey=E2=80=99s Cookies!=E2=80=9D</font></p></div></strong></font>\r\n\r\n<div><span color=3D\"#000000\" style=3D\"font-family: Arial, Helvetica, sans-=\r\nserif\" size=3D\"4;\"><strong>Go to&nbsp;</strong></span><strong><span  style=\r\n=3D\"font-family: Arial\" size=3D\"4;\">&nbsp;</span></strong></div>\r\n\r\n<div><strong><span  style=3D\"font-family: Arial\" size=3D\"4;\"><a href=3D\"ht=\r\ntp://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/h10\">www.tickfolio.com/a=\r\nrchives/1106</a> &nbsp;to&nbsp;order your Gourmet Cookies today!</span></s=\r\ntrong></div>\r\n\r\n<div><strong><span  style=3D\"font-family: Arial\" size=3D\"4;\">&nbsp;</span>=\r\n</strong></div>\r\n\r\n<div><strong><span  style=3D\"font-family: Arial;\">Comments or Questions?=\r\n <a href=3D\"mailto:info@Tickfolio.com\">info@Tickfolio.com</a></span></stro=\r\nng></div></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr><tr><td valign=3D\"top\" width=3D\"100%\"><img  style=3D\"display: bloc=\r\nk;\" alt=3D\"\" src=3D\"https://app.icontact.com/icp/loadimage.php/mogile/6159=\r\n54/b4bc5d5af51f975c01c81e2a038f6a25/image/jpeg\" complete=3D\"complete\"></td=\r\n><!-- br\r\n\r\n--></tr><tr><td valign=3D\"top\" width=3D\"100%\">\r\n\r\n<table width=3D\"700\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"0=\r\n\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"69\"><img  style=\r\n=3D\"display: block;\" alt=3D\"\" src=3D\"https://app.icontact.com/icp/loadimag=\r\ne.php/mogile/615954/fba9b691859d6cdfb483e09e57fda789/image/jpeg\" complete=\r\n=3D\"complete\" width=3D\"69\" height=3D\"279\"></td><!-- br\r\n\r\n--><td valign=3D\"top\" width=3D\"278\">\r\n\r\n<table width=3D\"278\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"2=\r\n0\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"100%\"><span st=\r\nyle=3D\"font-size: 28px;\" color=3D\"#000000\" style=3D\"font-family: Georgia,=\r\n Times New Roman, Times, serif\" size=3D\"5;\">Deal of the Week</span><br>\r\n\r\n<font color=3D\"#000000\" face=3D\"Arial, Helvetica, sans-serif\" size=3D\"2\">T=\r\nickfolio in partnership with Just Great Tickets offers outstanding deals=\r\n on Bears, Sox, Cubs, Ohio State and Nortre Dame Football tickets, Bulls,=\r\n Blackhawks and ALL Concerts and Theater for Individuals and Groups! (<str=\r\nong><a href=3D\"http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/h11\">Cl=\r\nick here</a></strong> to complete the inquiry form and receive a quote wit=\r\nhin 24 hours)<br>\r\n\r\n<br>\r\n\r\n<a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/S=\r\neMzKTgISMiH7cKVML-Qdg/h12\"><img  style=3D\"display: block;\" alt=3D\"\" src=3D=\r\n\"https://app.icontact.com/icp/loadimage.php/mogile/615954/27cade05fad83fd6=\r\n736d3907aab81341/image/jpeg\" complete=3D\"complete\" align=3D\"right\" border=\r\n=3D\"0\"></a></font></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--><td valign=3D\"top\" width=3D\"66\"><img  style=3D\"display: block;\" alt=3D\"=\r\n\" src=3D\"https://app.icontact.com/icp/loadimage.php/mogile/615954/fbf9f554=\r\nc74eb4d6b351237f5cfb2a6e/image/jpeg\" complete=3D\"complete\" width=3D\"66\" he=\r\night=3D\"286\"></td><!-- br\r\n\r\n--><td valign=3D\"top\" width=3D\"287\">\r\n\r\n<table width=3D\"287\" border=3D\"0\" bordercolor=3D\"#000000\" cellpadding=3D\"2=\r\n0\" cellspacing=3D\"0\"><tbody><tr><td valign=3D\"top\" width=3D\"100%\"><span st=\r\nyle=3D\"font-size: 28px;\" color=3D\"#000000\" style=3D\"font-family: Georgia,=\r\n Times New Roman, Times, serif\" size=3D\"5;\">JunoWallet</span><br>\r\n\r\n<font color=3D\"#000000\" face=3D\"Arial, Helvetica, sans-serif\" size=3D\"2\">T=\r\nickfolio in partnership with JunoWallet is pleased to offer FREE Promotion=\r\nal Mobile Gift Cards delivered by <strong><a href=3D\"http://u15167.sendgri=\r\nd.org/s/SeMzKTgISMiH7cKVML-Qdg/h13\">Tickfolio.com</a></strong>. No purchas=\r\ne is required to receive your free gift cards so enjoy our generous gift=\r\n to you our valued clients and prospects...<br>\r\n\r\n<br>\r\n\r\n<a style=3D\"text-decoration: none;\" href=3D\"http://u15167.sendgrid.org/s/S=\r\neMzKTgISMiH7cKVML-Qdg/h14\"><img  style=3D\"display: block;\" alt=3D\"\" src=3D=\r\n\"https://app.icontact.com/icp/loadimage.php/mogile/615954/27cade05fad83fd6=\r\n736d3907aab81341/image/jpeg\" complete=3D\"complete\" align=3D\"right\" border=\r\n=3D\"0\"></a></font></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr></tbody></table></td><!-- br\r\n\r\n--></tr></tbody></table><br>\r\n\r\n<span style=3D\"padding: 0px;\">&nbsp;</span></td><!-- br\r\n\r\n--></tr></tbody></table>If you'd like to unsubscribe and stop receiving th=\r\nese emails <a href=3D\"http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/=\r\nuh\">click here</a>.\r\n\r\n<img src=3D\"http://u15167.sendgrid.org/s/SeMzKTgISMiH7cKVML-Qdg/o0.gif\">",
				"text/plain":              "Content-Type: text/plain; charset=\"UTF-8\"\r\nContent-Disposition: inline\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\nTickfolio.com - Hot Deals\r\n\r\n\r\n\r\n$50 Value =E2=80=93 for $37.50 -25% Discount =E2=80=93 You Save $12.50\r\n\r\n\r\n\r\nGourmet Cookies Shipped Nationwide. Best Wishes, Thank You, Congrats, and=\r\n of course Happy Birthday... There are endless reasons to send your friend=\r\ns, family members, co-workers and clients they are truly amazing!\r\n\r\n\r\n\r\n-=3D Just Great Tickets =3D-\r\n\r\n\r\n\r\nTickfolio in partnership with Just Great Tickets offers outstanding deals=\r\n on Bears, Sox, Cubs, Ohio State and Nortre Dame Football tickets, Bulls,=\r\n Blackhawks and ALL Concerts and Theater for Individuals and Groups! (Go=\r\n to http://tickfolio.com/archives/1106?utm_source=3Dsendgrid.com&utm_mediu=\r\nm=3Demail&utm_campaign=3Dwebsite to complete the inquiry form and receive=\r\n a quote within 24 hours)\r\n\r\n\r\n\r\n-=3D JunoWallet =3D-\r\n\r\n\r\n\r\nTickfolio in partnership with JunoWallet is pleased to offer FREE Promotio=\r\nnal Mobile Gift Cards delivered by Tickfolio.com.\r\n\r\n\r\n\r\nVisit tickfolio.com for more details.If you'd like to unsubscribe and stop=\r\n receiving these emails click here: http://u15167.sendgrid.org/s/SeMzKTgIS=\r\nMiH7cKVML-Qdg/ut.",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.filename, func(t *testing.T) {
			mimeBytes, err := ioutil.ReadFile(tc.filename)
			assert.NoError(t, err)
			msg, err := Parse(string(mimeBytes))
			assert.NoError(t, err)
			defer msg.Close()

			actual := make(map[string]string)
			err = msg.Walk(func(p *Part) error {
				actual[p.ContentType()] = p.String()
				return nil
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPart_GetHeader(t *testing.T) {
	type TestCase struct {
		filename string
		expected []map[string]string
	}
	tcs := []TestCase{
		{
			filename: `test_data/mime_simple.eml`,
			expected: []map[string]string{
				{
					"Content-Type": "text/html",
				},
				{
					"Content-Transfer-Encoding": "quoted-printable",
					"Content-Type":              "text/plain; charset=us-ascii",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.filename, func(t *testing.T) {
			mimeBytes, err := os.ReadFile(tc.filename)
			assert.NoError(t, err)
			msg, err := Parse(string(mimeBytes))
			assert.NoError(t, err)
			defer msg.Close()

			section := 0
			err = msg.Walk(func(p *Part) error {
				actual := make(map[string]string)

				for k := range tc.expected[section] {
					actual[k] = p.GetHeader(k)
				}

				assert.Equal(t, tc.expected[section], actual)
				section += 1
				return nil
			})

			assert.NoError(t, err)
		})
	}
}

func TestPart_GetHeaders(t *testing.T) {
	type TestCase struct {
		filename string
		expected []map[string]string
	}
	tcs := []TestCase{
		{
			filename: `test_data/mime_simple.eml`,
			expected: []map[string]string{
				{
					"Content-Type": "text/html",
				},
				{
					"Content-Transfer-Encoding": "quoted-printable",
					"Content-Type":              "text/plain; charset=us-ascii",
					"X-Sg-Header":               "sg-header",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.filename, func(t *testing.T) {
			mimeBytes, err := os.ReadFile(tc.filename)
			assert.NoError(t, err)
			msg, err := Parse(string(mimeBytes))
			assert.NoError(t, err)
			defer msg.Close()

			actual := make([]map[string]string, len(tc.expected))
			section := 0
			err = msg.Walk(func(p *Part) error {
				actual[section] = make(map[string]string)
				if p.gmimePart == nil {
					return nil
				}

				// add headers into the actual map
				for k, v := range p.GetHeaders() {
					actual[section][k] = v[0]
				}

				section += 1
				return nil
			})
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

// BenchmarkPart_String benchmarks Part.String() method.
func BenchmarkPart_String(b *testing.B) {
	mimeBytes, _ := ioutil.ReadFile(`test_data/rfc822.eml`)
	msg, _ := Parse(string(mimeBytes))
	defer msg.Close()
	var part *Part
	_ = msg.Walk(func(p *Part) error {
		if p.ContentType() == `message/rfc822` {
			part = p
		}
		return nil
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		part.String()
	}
}

// TestPart_StringMemoryLeak tests for memory leak in this method
// This was created primarily for message/rfc822 format reading, so use that
// run this and watch the memory usage
// * go test -v . -test.run=TestPart_StringMemoryLeak
// * run activity monitor, find gmime.test process and make sure it doesn't jump up while it's running
func TestPart_StringMemoryLeak(t *testing.T) {
	t.Skipf("skipping since we don't need to run this normally, don't skip it if you want to test memory usage")
	mimeBytes, _ := ioutil.ReadFile(`test_data/rfc822.eml`)

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < 100000000; i++ {
		msg, _ := Parse(string(mimeBytes))
		var part *Part
		_ = msg.Walk(func(p *Part) error {
			if p.ContentType() == `message/rfc822` {
				part = p
			}
			return nil
		})
		part.String()
		msg.Close()
	}

	runtime.GC()
	runtime.ReadMemStats(&m2)

	fmt.Printf("total alloc: %d\n", m2.TotalAlloc-m1.TotalAlloc)
}
