#Not for production, only for testing purposes
FROM golang:alpine

WORKDIR /go/src/github.com/melchor629/speedy

RUN apk add --no-cache libpcap-dev git build-base && \
    go get -d -v github.com/influxdata/influxdb/client/v2 && \
    go get -d -v github.com/google/gopacket && \
    go get -d -v github.com/lib/pq

COPY . .
RUN go install -v ./...

CMD ["speedy"]