// Run
// cd go-gmime/bench
// go run gmime/test.go

package main

import (
	"bufio"
	"bytes"
	"github.com/sendgrid/go-gmime/bench/util"
	"github.com/sendgrid/go-gmime/gmime"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	root := "./data"
	filepath.Walk(root, util.BuildParsingVisitor(parseGMime))

	log.Printf("Parsed %d files in %f seconds with GMime binding.\n", util.Counter, util.Time)
}

var parseGMime util.ParseFunc = func(data []byte) bool {
	loops := 1

	if len(os.Args) > 1 {
		loops, _ = strconv.Atoi(os.Args[1])
	}

	for i := 0; i < loops; i++ {
		buf := bytes.NewBuffer(data)
		reader := bufio.NewReader(buf)
		parse := gmime.NewParse(reader)

		parse.Headers()
		parse.Header("Content-Type")
		text, _ := parse.Text()
		html, _ := parse.Html()
		attachments := parse.Attachment()

		if text == "" && html == "" && attachments == nil {
			return false
		}
	}

	return true
}
