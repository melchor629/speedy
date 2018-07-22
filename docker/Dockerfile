#Not for production, only for testing purposes
FROM golang:alpine as builder

WORKDIR /go/src/github.com/melchor629/speedy
COPY . .

RUN apk add --no-cache libpcap-dev git build-base && \
    go get -d -v github.com/influxdata/influxdb/client/v2 && \
    go get -d -v github.com/google/gopacket && \
    go get -d -v github.com/lib/pq && \
    go install -v ./...

FROM alpine:latest
COPY --from=builder /go/bin/* /usr/bin
RUN apk add --no-cache libpcap

CMD ["speedy"]