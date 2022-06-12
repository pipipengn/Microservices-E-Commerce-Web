FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o user-web .

WORKDIR /dist

RUN cp /build/user-web .

EXPOSE 8020

CMD ["/dist/user-web"]


