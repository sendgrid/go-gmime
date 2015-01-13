// Run
// cd go-gmime/bench
// go run native/test.go

package main

import (
	"github.com/le0pard/go-falcon/parser"
	"github.com/le0pard/go-falcon/protocol/smtpd"
	"github.com/sendgrid/go-gmime/bench/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// XXX We don't correctly newline-terminate our contentd samples,
// and that's why this parser chokes. Fix that problem with the
// following command (assuming STDERR was redirected to native.log):
// cat native.log | cut -d : -f 4 | sed "s/^ /.\//" | awk '{ printf "\n" >> $0; close $0 }'

func main() {
	root := "./data"
	filepath.Walk(root, util.BuildParsingVisitor(parseNative))

	log.Printf("Parsed %d files in %f seconds with Native implementation.\n", util.Counter, util.Time)
}

var parseNative util.ParseFunc = func(data []byte) bool {
	loops := 1

	if len(os.Args) > 1 {
		loops, _ = strconv.Atoi(os.Args[1])
	}

	for i := 0; i < loops; i++ {
		envelop := &smtpd.BasicEnvelope{MailboxID: 0, MailBody: data}
		myParser, err := parser.ParseMail(envelop)

		if err != nil {
			return false
		}

		headers := myParser.Headers
		headers.Get("Content-Type")
		text := myParser.TextPart
		html := myParser.HtmlPart
		attachments := myParser.Attachments

		// Just need to use the data above to make the compiler happy
		if text == "" && html == "" && attachments == nil {
			return true
		}
	}

	return true
}
