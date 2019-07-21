FROM arm64v8/golang:alpine as builder

RUN apk add --no-cache libpcap-dev git build-base && \
    go get -d -v github.com/influxdata/influxdb1-client/v2 && \
    go get -d -v github.com/google/gopacket && \
    go get -d -v github.com/lib/pq && \
    go get -d -v github.com/melchor629/speedy && \
    cd src/github.com/melchor629/speedy && \
    go install -v ./...

FROM arm64v8/alpine:slim
COPY --from=builder /go/bin/* /usr/bin
RUN apk add --no-cache libpcap

CMD ["speedy"]