# valgrind only seems to work correctly with Go <= 1.17 (seeing an unknown instruction
#   error on higher versions)
# admittedly, I haven't spent too much time trying to fix it but keep in mind
# if you're trying to increase the Go version
FROM golang:1.17-alpine3.16
RUN apk add gmime-dev build-base valgrind
COPY . /go/src/github.com/sendgrid/go-gmime
WORKDIR /go/src/github.com/sendgrid/go-gmime
RUN ["go", "build", "./cmd/gmime/main.go"]
