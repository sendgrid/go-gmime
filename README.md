go-gmime
========

This project is a binding of the C-based GMime library for Golang.

Dependencies
---

- Go v. 1.2+
- GMime v. 2.6.19 or later, built with glib 2.38.0 or later.

Usage
---

Just add library with:

    go get github.com/sendgrid/go-gmime

or using your favorite dependency manager (for example `godep` or `gom`)

Setup for development
---

For OSX

	# install required packages
	brew install go
	brew install gmime
	bin/install_deps

On Ubuntu/Debian (trusty or jessie/sid, or later):

    # install go and required packages
    sudo apt-get install golang libgmime-2.6-dev

Earlier Ubuntu like `12.04 (precise)` require to backport glib2 (version
2.38.1 or later) and gmime-2.6 (2.6.19 or later)
Follow instructions in ``doc/packages.md`` for details.

``bin/install_dep`` also put symlink to package dir, so cloned dir accessable
for import by full-qualified name ``github.com/sendgrid/go-gmime/gmime``

Compose Message
---

```go

composer := NewComposer()
composer.AddFrom("Good Sender <good_sender@example.com>")
composer.AddTo("good_customer@example.com", "Good Customer")

// read data from a file:
fileHandler, _ := os.Open("data.txt")
defer fileHandler.Close()
reader := bufio.NewReader(fileHandler)

composer.AddHTML(reader)
print(composer.GetMessage())
```

Parse Message
---

```go

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

```

Documentation
---

Use [https://developer.gnome.org/gmime/stable/|GMime Reference Manual] as
documentation, this library have all classes/methods named very similiar with
original ``GMime`` classes. 

All objects will be unreferenced (in terms of glib/gmime) when go-gmime
provided wrappers will be collected by go runtime.


Testing / Coverage
---

	# run all tests on host machine
	bin/test
	
	# run all tests on Vagrant CentOS VM
	bin/testvm
	
	# run a specific test:
	bin/test TestSimpleMessage
		
	# generate tests coverage
	bin/cover

Vagrant VMs 
---
	We use Ubuntu to run Valgrind and CentOS to match SendGrid production
	
	# setup the Vagrant vm box and access it:
	vagrant up go-gmime-centos
	vagrant up go-gmime-ubuntu
	
	# access a VM box
	vagrant ssh go-gmime-centos
	vagrant ssh go-gmime-ubuntu
	
	# run test on VM
	bin/testvm

Benchmarks
---
Create and copy test emails to the data/ directory at the root of this project. 

	# benchmark using GMime binding
	go run bench/gmime/test.go

	# benchmark using native Go
	go run bench/native/test.go
	
	# benchmark using Perl
	perl bench/perl/test.pl
	
	# run integration GMime binding benchmark
	bin/bench

Memory Check
---
We use Valgrind with **go-gmime-ubuntu** VM to check for memory leak. Run the command below inside the VM and check valgrind.log

	bin/valgrind	

(it require valgrind 3.9.x, 3.10.x is better, and gcc-go at least 4.9)

Contributing
---
	# Don't forget to run gofmt before commiting
	bin/fmt

License
---

go-gmime is licensed under the terms of the MIT license reproduced below.
This means that go-gmime is free software and can be used for both academic
and commercial purposes at absolutely no cost.
