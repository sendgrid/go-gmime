package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/sendgrid/go-gmime/v1/gmime"
)

func ExampleReadingHeaders() {
	fileHandle, err := os.Open("gmime/fixtures/text_attachment.eml")
	if err != nil {
		panic(err)
	}
	defer fileHandle.Close()

	reader := bufio.NewReader(fileHandle)
	data, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	msg, err := gmime.Parse(string(data))
	if err != nil {
		panic(err)
	}
	defer msg.Close()

	// How to loop over headers:
	for headerName, headerValues := range msg.Headers() {
		fmt.Printf("%s: %v\n", headerName, headerValues)
	}

	// Setting subject line
	msg.SetSubject("My Favorite Subject!")
	v, err := msg.Export()
	if err != nil {
		panic(err)
	}

	fmt.Println("****** Rendered MIME******")
	fmt.Println(string(v))
}

func main() {
	ExampleReadingHeaders()
	debug.FreeOSMemory()
}
