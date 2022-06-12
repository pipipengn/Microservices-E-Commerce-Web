FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o user-srv .

WORKDIR /dist

RUN cp /build/user-srv .

EXPOSE 50051

CMD ["/dist/user-srv"]


