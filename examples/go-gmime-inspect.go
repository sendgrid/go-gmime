// This file is brainless porting of examples/imap-example.go
// This code haven't anything common or related with IMAP, I haven't any idea
// why gmime developers used this name.

package main

import (
	"flag"
	"fmt"
	"github.com/sendgrid/go-gmime/gmime"
)

var file_name string = ""
var scan_from bool

func init() {
	flag.StringVar(&file_name, "filename", "file.eml", "email file to parse")
	flag.BoolVar(&scan_from, "from", false, "scan from")
}

func main() {
	flag.Parse()
	if file_name == "" {
		panic("no filename given")
	}

	fs := gmime.NewFileStreamForPath(file_name, "r")
	if fs == nil {
		panic("can't open " + file_name)
	}

	parser := gmime.NewParserWithStream(fs)
	parser.SetScanFrom(scan_from)
	message := parser.ConstructMessage()
	if message != nil {
		msgid, ok := message.MessageId()
		fmt.Printf("Inspecting %s\n: message.MessageId() == \"%v\", %v\n", file_name, msgid, ok)
		inspect_message(message)
	}
}

func inspect_message(m gmime.Message) {
	inspect_header(m)
	inspect_part(m.MimePart(), 0)
}

func inspect_header(m gmime.Message) {
	fmt.Printf("=================== HEADER =======================\n")
	fmt.Printf("%v\n", m.Headers())
	fmt.Printf("=================== END HEADER ===================\n")
}

func inspect_part(ob gmime.Object, level int) {
	{
		fmt.Printf("=================== PART HEADER =======================\n")
		fmt.Printf("%v\n", ob.Headers())
		fmt.Printf("================= END PART HEADER =====================\n")
	}

	fmt.Printf("=================== PART BODY =======================\n")
	if mp, ok := ob.(gmime.Multipart); ok {
		n := mp.Count()
		fmt.Printf("Object type: gmime.Multipart, %d items\n", n)
		for i := 0; i < n; i++ {
			subpart := mp.GetPart(i)
			fmt.Printf("%d >>> SUBPART: %d\n", level, i)
			inspect_part(subpart, level+1)
		}
	} else if mp, ok := ob.(gmime.MessagePart); ok {
		fmt.Printf("Object type: gmime.MessagePart\n")
		ostream := gmime.NewMemStream()
		mp.Message().WriteToStream(ostream)
		defer inspectAndClose(ostream)
	} else if p, ok := ob.(gmime.Part); ok {
		fmt.Printf("Object type: gmime.Part\n")
		ostream := gmime.NewMemStream()
		defer inspectAndClose(ostream)
		dataWrapper := p.ContentObject()
		dataWrapper.WriteToStream(ostream)
	}
	fmt.Printf("=================== END OF PART BODY =======================\n")
}

func inspectAndClose(ostream gmime.MemStream) {
	defer ostream.Close()
	fmt.Printf("=================== MEM STREAM =======================\n")
	fmt.Printf("%v", string(ostream.Bytes()))
	fmt.Printf("=================== END OF MEM STREAM =======================\n")
}
