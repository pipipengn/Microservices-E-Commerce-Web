FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o order-web .

WORKDIR /dist

RUN cp /build/order-web .

EXPOSE 8022

CMD ["/dist/order-web"]

