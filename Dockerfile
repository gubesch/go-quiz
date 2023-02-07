# Build
FROM golang:1.20.0-alpine3.17 AS base
COPY . /go/src/go-quiz

WORKDIR /go/src/go-quiz
RUN apk add git gcc

# Go modules
ENV GO111MODULE=on
#RUN go mod tidy
RUN go mod download

# Compile
RUN go build -a -tags netgo -ldflags '-w' -o /go/bin/go-quiz /go/src/go-quiz/main.go

# Package
FROM alpine:3.17
COPY --from=base /go/bin/go-quiz /go-quiz/go-quiz
COPY --from=base /go/src/go-quiz/migration /go-quiz/migration
WORKDIR /go-quiz
ENTRYPOINT ["/bin/sh"]
