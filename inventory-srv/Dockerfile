FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o inventory-srv .

WORKDIR /dist

RUN cp /build/inventory-srv .

EXPOSE 50053

CMD ["/dist/inventory-srv"]

