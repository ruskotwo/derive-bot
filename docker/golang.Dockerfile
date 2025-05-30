# Building
FROM golang:1.23-alpine as builder

#[For Dev]
#RUN go install github.com/go-delve/delve/cmd/dlv@v1.23.1

RUN apk update && \
    apk add --no-cache \
    git \
    openssh \
    gcc \
    libc-dev \
    ca-certificates

ADD . /go/src/github.com/ruskotwo/derive-bot
WORKDIR /go/src/github.com/ruskotwo/derive-bot

RUN go mod download -x

#[For Dev]
#RUN go build -v -race -gcflags "all=-N -l" -o /go/bin/derive-bot ./main.go
#[For All]
RUN go build -v  -o /go/bin/derive-bot ./main.go

#Running
FROM alpine:3.20

COPY --from=builder /go/bin/derive-bot /usr/local/bin/derive-bot

COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

#[For Dev]
#COPY --from=builder /go/bin/dlv /go/bin/dlv

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
