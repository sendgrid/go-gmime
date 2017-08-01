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
	msg := gmime.Parse(string(data))
	defer msg.Close()
	println("Envelope Subject: ", msg.Subject())
	println("Envelope Content-Type:", msg.ContentType())
	println("Envelope Message-ID", msg.Header("Message-ID"))
	fmt.Println("All Headers:", msg.Headers())

	msg.Walk(func(p *gmime.Part) {
		println("content-type:", p.ContentType())
		if p.IsText() {
			println("text:", p.Text())
			p.SetText("my replaced всякий текст スラングまで幅広く収録")
		} else {
			// fmt.Println("Bytes:", string(p.Bytes()))
		}
	})
	println(">>> test")

	// data, _ = msg.Export()
	// println(">>>>>>>>>>>>> Export: ", string(data))
}
