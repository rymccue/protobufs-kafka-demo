FROM golang:1.10-alpine

LABEL authors="Ryan McCue <ryan@msys.ca>"

RUN apk add --no-cache ca-certificates openssl git
RUN wget -O /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
  chmod +x /usr/local/bin/dep

RUN mkdir /go/src/app

ADD . /go/src/app/

WORKDIR /go/src/app

run go get github.com/rymccue/protobufs-kafka-demo

RUN go get github.com/golang/protobuf/proto

RUN go get github.com/Shopify/sarama

RUN go build -o main .

CMD ["/go/src/app/main"]