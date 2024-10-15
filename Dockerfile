
FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./. ./

RUN go build main.go

FROM busybox

WORKDIR /go

COPY --from=builder /go/src/main /go/main

EXPOSE 8080

ENTRYPOINT [ "/go/main" ]
