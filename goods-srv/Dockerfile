FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o goods-srv .

WORKDIR /dist

RUN cp /build/goods-srv .

EXPOSE 50052

CMD ["/dist/goods-srv"]

