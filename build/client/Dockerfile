FROM golang:alpine as builder
WORKDIR /go/src/github.com/xdefrag/hash-ketchum
RUN apk --no-cache add git make
ENV GO111MODULE=on
COPY . ./
RUN make build-client

FROM scratch
COPY --from=builder /go/src/github.com/xdefrag/hash-ketchum/dist/client .
CMD ["./client"]
