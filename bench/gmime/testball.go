// Run
// cd go-gmime/bench
// go run gmime/test.go

package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/sendgrid/go-gmime/bench/util"
	"github.com/sendgrid/go-gmime/gmime"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"
)

var archive string
var loops int
var quiet bool

func init() {
	flag.StringVar(&archive, "archive", "contentd.tar.gz", ".tar.gz archive to parse")
	flag.IntVar(&loops, "l", 1, "How many times repeat parsing")
	flag.BoolVar(&quiet, "q", false, "report only in end of cycle")
}

func main() {
	flag.Parse()
	for i := 0; i < loops; i++ {
		parseArchive(archive)
	}
	freeOSMemory()
	gmime.Shutdown()
	freeOSMemory()
}

func freeOSMemory() {
	for i := 0; i < 10; i++ {
		debug.FreeOSMemory()
	}
}

func parseArchive(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic("can't open " + filename)
	}
	defer f.Close()

	gz, err2 := gzip.NewReader(f)
	if err2 != nil {
		panic("bad gzip header")
	}
	defer gz.Close()

	t := tar.NewReader(gz)

	for {
		hdr, err := t.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		path := hdr.Name
		start := time.Now().UnixNano()
		if parseGMime(t) {
			end := time.Now().UnixNano()
			seconds := float64(end-start) / 1000000000.0
			util.Time += seconds
			util.Counter++
			if !quiet {
				fmt.Printf("%s,%f\n", path, seconds)
			}
		} else {
			log.Printf("Failed to parse: %s\n", path)
		}
	}

	log.Printf("Parsed %d files in %f seconds with GMime binding.\n", util.Counter, util.Time)
}

func parseGMime(reader io.Reader) bool {
	parse := gmime.NewParse(reader)

	parse.Headers()
	parse.Header("Content-Type")
	text := parse.Text()
	html := parse.Html()
	attachments := parse.Attachment()

	if text == "" && html == "" && attachments == nil {
		return false
	}

	return true
}
