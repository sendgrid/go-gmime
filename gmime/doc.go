/*
	Package gmime provides the binding to the C-based GMime

	Compose Message
    ---

	composer := NewComposer()
	composer.AddFrom("Good Sender <good_sender@example.com>")
	composer.AddTo("good_customer@example.com", "Good Customer")

	// read data from a file:
	fileHandler, _ := os.Open("data.txt")
	defer fileHandler.Close()
	reader := bufio.NewReader(fileHandler)

	composer.AddHTML(reader)
	print(composer.GetMessage())


	Parse Message
	---

	// Example on how to create a NewParse

	fileHandler, _ := os.Open("fixtures/text_attachment.eml")
	defer fileHandler.Close()
	reader := bufio.NewReader(fileHandler)
	parse := NewParse(reader)

	fmt.Println(parse.From())
	fmt.Println(parse.To())
	fmt.Println(parse.Subject())

	// Output:
	// Dave McGuire <david.mcguire@sendgrid.com>
	// Team R&D <team-rd@sendgrid.com>
	// Distilled Content-Type headers (outbound & failed)
*/

package gmime
