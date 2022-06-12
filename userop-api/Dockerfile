FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .

RUN go build -o userop-web .

WORKDIR /dist

RUN cp /build/userop-web .

EXPOSE 8023

CMD ["/dist/userop-web"]

