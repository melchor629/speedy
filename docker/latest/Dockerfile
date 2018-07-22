FROM golang:latest as builder

WORKDIR /go/src/speedy
COPY . .

RUN go get -d -v github.com/influxdata/influxdb/client/v2 && \
    go get -d -v github.com/google/gopacket && \
    apt-get update && \
    apt-get install --no-install-recommends -y libpcap-dev && \
    go install -v ./...

FROM debian:latest
COPY --from=builder /go/bin/* /usr/bin
RUN apt-get update && \
    apt-get install --no-install-recommends -y libpcap0.8 && \
    rm -r /var/lib/apt && \
    rm -r /var/cache/apt

CMD ["speedy"]