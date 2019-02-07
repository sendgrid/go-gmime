package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sendgrid/go-gmime/gmime"
)

func main() {
	println(">>> test")

	fh, err := os.Open("./gmime/fixtures/inline-attachment.eml")
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	reader := bufio.NewReader(fh)
	data, _ := ioutil.ReadAll(reader)
	msg, err := gmime.Parse(string(data))
	if err != nil {
		panic(err)
	}
	defer msg.Close()

	// test stuff goes here

	b, err := msg.Export()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	fmt.Printf("%s\n", string(b))

	// data, _ = msg.Export()
	// println(">>>>>>>>>>>>> Export: ", string(data))
}
