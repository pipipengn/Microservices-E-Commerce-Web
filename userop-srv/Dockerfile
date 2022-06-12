FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o userop-srv .

WORKDIR /dist

RUN cp /build/userop-srv .

EXPOSE 50055

CMD ["/dist/userop-srv"]

