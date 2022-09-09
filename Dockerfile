FROM golang:1.19-alpine3.16 as builder
RUN apk add gmime-dev build-base
COPY . /go/src/github.com/sendgrid/go-gmime
RUN ["go", "build", "/go/src/github.com/sendgrid/go-gmime/gmime"]
