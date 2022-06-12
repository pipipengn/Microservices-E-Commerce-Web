FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o order-srv .

WORKDIR /dist

RUN cp /build/order-srv .

EXPOSE 50054

CMD ["/dist/order-srv"]

