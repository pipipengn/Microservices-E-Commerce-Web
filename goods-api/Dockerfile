FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o goods-web .

WORKDIR /dist

RUN cp /build/goods-web .

EXPOSE 8021


CMD ["/dist/goods-web"]

