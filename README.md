go-gmime
========

This project is a binding of the C-based GMime library for Golang.

Dependencies
---

- Go1.17
- GMime >= v3.0 

Usage
---

Just add library with:

    go get github.com/sendgrid/go-gmime/gmime

Setup for development
---

For OSX

	# install required packages
	brew install go gmime

Examples
---
See cmd/gmime/main.go for examples

Testing / Coverage
---

	# run all tests on host machine
    go test ./gmime/...
	
	# run all tests in Docker container
    docker build . -t docker.io/library/go-gmime
    docker run -it -v $(pwd):/go/src/github.com/sendgrid/go-gmime docker.io/library/go-gmime
    go test ./gmime/...

Memory Check
---
We use Valgrind to check for memory leaks. The provided Dockerfile should setup a container ready to run Valgrind

	bin/valgrind	

Contributing
---
	# Don't forget to run gofmt before commiting
	go fmt ./...

Afterwards, submit a PR for review.

License
---

go-gmime is licensed under the terms of the MIT license reproduced below.
This means that go-gmime is free software and can be used for both academic
and commercial purposes at absolutely no cost.
