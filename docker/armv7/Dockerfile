FROM arm32v7/golang:alpine as builder

RUN apk add --no-cache libpcap-dev git build-base && \
    go get -d -v github.com/influxdata/influxdb1-client/v2 && \
    go get -d -v github.com/google/gopacket && \
    go get -d -v github.com/lib/pq && \
    go get -d -v github.com/melchor629/speedy && \
    echo "do" && \
    cd src/github.com/melchor629/speedy && \
    echo "nothing" && \
    go install -v ./...

FROM arm32v7/alpine
COPY --from=builder /go/bin/* /usr/bin
RUN apt-get update && \
    apt-get install --no-install-recommends -y libpcap0.8 && \
    rm -r /var/lib/apt && \
    rm -r /var/cache/apt

CMD ["speedy"]
