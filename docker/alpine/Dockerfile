FROM golang:alpine as builder

WORKDIR /go/src/speedy
COPY . .

RUN apk add --no-cache libpcap-dev git build-base && \
    go get -d -v github.com/influxdata/influxdb/client/v2 && \
    go get -d -v github.com/google/gopacket && \
    go install -v ./...

FROM alpine:latest
COPY --from=builder /go/bin/* /usr/bin
RUN apk add --no-cache libpcap

CMD ["speedy"]