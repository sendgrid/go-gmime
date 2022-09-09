// Run
// cd go-gmime/bench
// go run gmime/test.go

package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sendgrid/go-gmime/v1/bench/util"
	"github.com/sendgrid/go-gmime/v1/gmime"
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
		text := parse.Text()
		html := parse.Html()
		attachments := parse.Attachment()

		if text == "" && html == "" && attachments == nil {
			return false
		}
	}

	return true
}
